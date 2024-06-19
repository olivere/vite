package vite

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
)

// Handler serves files from the Vite output directory.
type Handler struct {
	fs              fs.FS
	fsFS            http.FileSystem
	fsHandler       http.Handler
	pub             fs.FS
	pubFS           http.FileSystem
	pubHandler      http.Handler
	manifest        *Manifest
	isDev           bool
	viteURL         string
	templates       map[string]*template.Template
	defaultMetadata *Metadata
}

// Config is the configuration for the handler.
type Config struct {
	// FS is the file system to serve files from. In production, this is
	// the Vite output directory, which usually is the "dist" directory.
	// In development, this is usually the root directory of the Vite app.
	FS fs.FS
	// PublicFS is the file system to serve public files from. This is
	// usually the "public" directory. It is optional and can be nil.
	// If it is nil, we will check if the "public" directory exists in
	// the Vite app, and serve files from there. If it does not exist,
	// we will not serve any public files. It is only used in development
	// mode.
	PublicFS fs.FS
	// IsDev is true if the server is running in development mode, false
	// otherwise.
	IsDev bool
	// ViteURL is the URL of the Vite server, used to load the Vite client
	// in development mode. It is unused in production mode.
	ViteURL string
}

// NewHandler creates a new handler.
//
// fs is the file system to serve files from, the Vite output directory
// (which usually is the "dist" directory). isDev is true if the server is
// running in development mode, false otherwise. viteServer is the URL of the
// Vite server, used to load the Vite client in development mode.
func NewHandler(config Config) (*Handler, error) {
	if config.FS == nil {
		return nil, fmt.Errorf("vite: fs is nil")
	}
	h := &Handler{
		fs:        config.FS,
		fsFS:      http.FS(config.FS),
		fsHandler: http.FileServerFS(config.FS),
		isDev:     config.IsDev,
		viteURL:   config.ViteURL,
		templates: make(map[string]*template.Template),
	}

	h.templates["index.html"] = template.Must(template.New("index.html").Parse(indexHTML))

	if !h.isDev {
		// We expect the output directory to contain a .vite/manifest.json file.
		// This file contains the mapping of the original file paths to the
		// transformed file paths.
		mf, err := h.fs.Open(".vite/manifest.json")
		if err != nil {
			return nil, fmt.Errorf("vite: open manifest: %w", err)
		}
		defer mf.Close()

		// Read the manifest file.
		h.manifest, err = ParseManifest(mf)
		if err != nil {
			return nil, fmt.Errorf("vite: parse manifest: %w", err)
		}
	} else {
		if h.viteURL == "" {
			h.viteURL = "http://localhost:5173"
		}

		if config.PublicFS == nil {
			// We will peek into the "public" directory of the Vite app, and
			// serve files from there (if it exists).
			pub, err := fs.Sub(config.FS, "public")
			if err == nil {
				h.pub = pub
				h.pubFS = http.FS(h.pub)
				h.pubHandler = http.FileServerFS(h.pub)
			}
		} else {
			h.pub = config.PublicFS
			h.pubFS = http.FS(config.PublicFS)
			h.pubHandler = http.FileServerFS(config.PublicFS)
		}
	}

	return h, nil
}

// SetDefaultMetadata sets the default metadata to use when rendering the
// page. This metadata is used when the context does not have any metadata.
func (h *Handler) SetDefaultMetadata(md *Metadata) {
	h.defaultMetadata = md
}

// RegisterTemplate registers a template with the handler. Notice that the
// template name must be unique, and "index.html" is already registered by
// default. If the template name is already registered, it panics.
func (h *Handler) RegisterTemplate(name, text string) {
	if h.templates == nil {
		h.templates = make(map[string]*template.Template)
	}
	if _, ok := h.templates[name]; ok {
		panic(fmt.Sprintf("vite: template %q already registered", name))
	}
	h.templates[name] = template.Must(template.New(name).Parse(text))
}

// HandlerFunc returns a http.HandlerFunc for h.
func (h *Handler) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(h.ServeHTTP)
}

// ServeHTTP handles HTTP requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Get the path from the URL: https://localhost/articles/123 -> /articles/123
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusBadRequest)
		return
	}

	isIndexPath := path == "/" || path == "/index.html"

	// Check if the file exists in the public directory.
	if h.isDev && h.pubFS != nil && h.pubHandler != nil && !isIndexPath {
		if _, err := h.pubFS.Open(path); err == nil {
			h.pubHandler.ServeHTTP(w, r)
			return
		}
	}

	// If we can, we serve from the file system
	_, err = h.fsFS.Open(path)
	if err != nil || isIndexPath {
		// We didn't find it in the file system, so we generate the HTML
		// from the entry point with Go templating.
		h.renderPage(w, r, path, nil)
		return
	}

	// Serve the file using the file server.
	h.fsHandler.ServeHTTP(w, r)
}

// pageData is passed to the template when rendering the page.
type pageData struct {
	IsDev               bool
	ViteURL             string
	Metadata            template.HTML
	PluginReactPreamble template.HTML
	StyleSheets         template.HTML
	Modules             template.HTML
	PreloadModules      template.HTML
	Scripts             template.HTML
}

// renderPage renders the page using the template.
func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, path string, chunk *Chunk) {
	page := pageData{
		IsDev:   h.isDev,
		ViteURL: h.viteURL,
	}

	// Inject metadata into the page.
	ctx := r.Context()
	md := MetadataFromContext(ctx)
	if md == nil {
		md = h.defaultMetadata
	}
	if md != nil {
		page.Metadata = template.HTML(md.String())
	}

	// Inject scripts into the page.
	scripts := ScriptsFromContext(ctx)
	if scripts != "" {
		page.Scripts = template.HTML(scripts)
	}

	// Handle both development and production modes.
	if h.isDev {
		page.PluginReactPreamble = template.HTML(PluginReactPreamble(h.viteURL))
	} else {
		if chunk == nil {
			chunk = h.manifest.GetEntryPoint()
			if chunk == nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		page.StyleSheets = template.HTML(h.manifest.GenerateCSS(chunk.Src))
		page.Modules = template.HTML(h.manifest.GenerateModules(chunk.Src))
		page.PreloadModules = template.HTML(h.manifest.GeneratePreloadModules(chunk.Src))
	}

	// Find the template to use.
	tmpl, ok := h.templates[path]
	if !ok {
		// Use index.html by default
		tmpl = h.templates["index.html"]
	}

	// Execute the template.
	if err := tmpl.Execute(w, page); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

var (
	indexHTML = `<!doctype html>
<html lang="en" class="h-full scroll-smooth">
  <head>
    <meta charset="UTF-8" />
	{{- if .Metadata }}
		{{ .Metadata }}
	{{- end }}
	{{- if .IsDev }}
		{{ .PluginReactPreamble }}
		<script type="module" src="{{ .ViteURL }}/@vite/client"></script>
		<script type="module" src="{{ .ViteURL }}/src/main.tsx"></script>
	{{- else }}
		{{- if .StyleSheets }}
		{{ .StyleSheets }}
		{{- end }}
		{{- if .Modules }}
		{{ .Modules }}
		{{- end }}
		{{- if .PreloadModules }}
		{{ .PreloadModules }}
		{{- end }}
	{{- end }}
	{{- if .Scripts }}
		{{ .Scripts }}
	{{- end }}
 </head>
  <body class="min-h-screen antialiased">
    <div id="root"></div>
  </body>
</html>
`
)
