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
	hfs             http.FileSystem
	fileServer      http.Handler
	manifest        *Manifest
	isDev           bool
	viteServer      string
	templates       map[string]*template.Template
	defaultMetadata *Metadata
}

// NewHandler creates a new handler.
//
// fs is the file system to serve files from, the Vite output directory
// (which usually is the "dist" directory). isDev is true if the server is
// running in development mode, false otherwise. viteServer is the URL of the
// Vite server, used to load the Vite client in development mode.
func NewHandler(fs fs.FS, isDev bool, viteServer string) (*Handler, error) {
	h := &Handler{
		fs:         fs,
		hfs:        http.FS(fs),
		isDev:      isDev,
		viteServer: viteServer,
		templates:  make(map[string]*template.Template),
	}
	h.fileServer = http.FileServer(h.hfs)

	h.templates["index.html"] = template.Must(template.New("index.html").Parse(indexHTML))

	if !isDev {
		// We expect the output directory to contain a .vite/manifest.json file.
		// This file contains the mapping of the original file paths to the
		// transformed file paths.
		mf, err := fs.Open(".vite/manifest.json")
		if err != nil {
			return nil, fmt.Errorf("vite: open manifest: %w", err)
		}
		defer mf.Close()

		// Read the manifest file.
		h.manifest, err = ParseManifest(mf)
		if err != nil {
			return nil, fmt.Errorf("vite: parse manifest: %w", err)
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

	// If we can, we serve from the file system
	_, err = h.hfs.Open(path)
	if err != nil || path == "/" || path == "/index.html" {
		// We didn't find it in the file system, so we generate the HTML
		// from the entry point with Go templating.
		h.renderPage(w, r, path, nil)
		return
	}

	// Serve the file using the file server.
	h.fileServer.ServeHTTP(w, r)
}

// pageData is passed to the template when rendering the page.
type pageData struct {
	IsDev               bool
	ViteDevServer       string
	Metadata            template.HTML
	PluginReactPreamble template.HTML
	StyleSheets         template.HTML
	Modules             template.HTML
	PreloadModules      template.HTML
}

// renderPage renders the page using the template.
func (h *Handler) renderPage(w http.ResponseWriter, r *http.Request, path string, chunk *Chunk) {
	page := pageData{
		IsDev:         h.isDev,
		ViteDevServer: h.viteServer,
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

	// Handle both development and production modes.
	if h.isDev {
		page.PluginReactPreamble = template.HTML(PluginReactPreamble(h.viteServer))
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
		<script type="module" src="{{ .ViteDevServer }}/@vite/client"></script>
		<script type="module" src="{{ .ViteDevServer }}/src/main.tsx"></script>
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
 </head>
  <body class="min-h-screen antialiased">
    <div id="root"></div>
  </body>
</html>
`
)
