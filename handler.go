package vite

import (
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
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
	viteEntry       string
	viteURL         string
	viteTemplate    Scaffolding
	templates       map[string]*template.Template
	defaultMetadata *Metadata
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
		fs:           config.FS,
		fsFS:         http.FS(config.FS),
		fsHandler:    http.FileServerFS(config.FS),
		isDev:        config.IsDev,
		viteEntry:    config.ViteEntry,
		viteURL:      config.ViteURL,
		viteTemplate: config.ViteTemplate,
		templates:    make(map[string]*template.Template),
	}

	// We register a fallback template.
	h.templates[fallbackTemplateName] = template.Must(template.New(fallbackTemplateName).Parse(fallbackHTML))

	if !h.isDev {
		// Production mode.
		//
		// We expect the output directory to contain a .vite/manifest.json file.
		// This file contains the mapping of the original file paths to the
		// transformed file paths.
		if config.ViteManifest == "" {
			config.ViteManifest = ".vite/manifest.json"
		}
		mf, err := h.fs.Open(config.ViteManifest)
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
		// Development mode.
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

// RegisterTemplate adds a new template to the handler's template collection.
// The 'name' parameter should match the URL path where the template will be used.
// Use "index.html" for the root URL ("/").
//
// Parameters:
//   - name: String identifier for the template, corresponding to its URL path
//   - text: String content of the template
//
// Panics if a template with the given name is already registered.
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
	// Normalize the path, e.g. /..//articles/123/ -> /articles/123
	path := path.Clean(r.URL.Path)

	isIndexPath := path == "/" || path == "/index.html"

	// Check if the file exists in the public directory.
	if h.isDev && h.pubFS != nil && h.pubHandler != nil && !isIndexPath {
		if _, err := h.pubFS.Open(path); err == nil {
			h.pubHandler.ServeHTTP(w, r)
			return
		}
	}

	if isIndexPath {
		// We didn't find it in the file system, so we generate the HTML
		// from the entry point with Go templating.
		h.renderPage(w, r, path, nil)
		return
	}

	if _, ok := h.templates[path]; ok {
		// We found a template for the path, so we render the page using
		// the template.
		h.renderPage(w, r, path, nil)
		return
	}

	// Check if the file exists in the file system.
	if _, err := h.fsFS.Open(path); err != nil {
		// The file does not exist in the file system, so 404.
		http.NotFound(w, r)
		return
	}

	// Serve the file using the file server.
	h.fsHandler.ServeHTTP(w, r)
}

// pageData is passed to the template when rendering the page.
type pageData struct {
	IsDev               bool
	ViteEntry           string
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
		IsDev:     h.isDev,
		ViteEntry: h.viteEntry,
		ViteURL:   h.viteURL,
	}

	// Inject metadata in// Check if the specified Vite template requires a preamble and set the
	// corresponding preamble string in the plugin configuration.
	//
	// If the Vite template value is less than 1, it is considered as an
	// uninitialized state, and the default React preamble is applied.
	// Otherwise, if the template requires a preamble, it uses the
	// specific preamble for the given Vite template.to the page.
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
		// Check if the specified Vite template requires a preamble and set the
		// corresponding preamble string in the plugin configuration.
		//
		// If the Vite template value is less than 1, it is considered as an
		// uninitialized state, and the default React preamble is applied.
		// Otherwise, if the template requires a preamble, it uses the
		// specific preamble for the given Vite template.
		if h.viteTemplate < 1 {
			page.PluginReactPreamble = template.HTML(React.Preamble(h.viteURL))
		} else if h.viteTemplate.RequiresPreamble() {
			page.PluginReactPreamble = template.HTML(h.viteTemplate.Preamble(h.viteURL))
		}
		// page.PluginReactPreamble = template.HTML(PluginReactPreamble(h.viteURL))
	} else {
		if chunk == nil {
			if page.ViteEntry == "" {
				chunk = h.manifest.GetEntryPoint()
			} else {
				entries := h.manifest.GetEntryPoints()
				for _, entry := range entries {
					if page.ViteEntry == entry.Src {
						chunk = entry
						break
					}
				}
			}
			if chunk == nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
		}
		assetsPrefix := ""
		page.StyleSheets = template.HTML(h.manifest.GenerateCSS(chunk.Src, assetsPrefix))
		page.Modules = template.HTML(h.manifest.GenerateModules(chunk.Src, assetsPrefix))
		page.PreloadModules = template.HTML(h.manifest.GeneratePreloadModules(chunk.Src, assetsPrefix))
	}

	var tmplName string
	if path == "/" {
		tmplName = "index.html"
	} else {
		tmplName = path
	}

	// Find the template by name.
	tmpl, ok := h.templates[tmplName]

	// Catch common variations. If a template isn't found by the exact name,
	// check for variations like: "page", "page.html", or "/page.html", to match
	// how users might have registered the template.
	if !ok {
		variations := []string{
			strings.TrimPrefix(tmplName, "/"),
			strings.TrimPrefix(tmplName, "/") + ".html",
			strings.TrimSuffix(strings.TrimPrefix(tmplName, "/"), ".html"),
			tmplName + ".html",
		}
		for _, variant := range variations {
			if t, found := h.templates[variant]; found {
				tmpl = t
				ok = true
				break
			}
		}
	}

	// Handle case when requested template is not found:
	// 1. If multiple templates exist, log a warning with the requested and available templates.
	// 2. Fall back to a default template.
	if !ok {
		if len(h.templates) > 1 {
			keys := make([]string, 0, len(h.templates))
			for k := range h.templates {
				keys = append(keys, k)
			}
			slog.Warn(
				"Template not found",
				"requestedTemplate", tmplName,
				"availableTemplates", strings.Join(keys, ", "),
			)
		}
		tmpl = h.templates[fallbackTemplateName]
	}

	// Execute the template.
	if err := tmpl.Execute(w, page); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

const fallbackTemplateName = "fallback.html"

var (
	fallbackHTML = `<!doctype html>
<html lang="en" class="h-full scroll-smooth">
  <head>
    <meta charset="UTF-8" />
	{{- if .Metadata }}
		{{ .Metadata }}
	{{- end }}
	{{- if .IsDev }}
		{{ .PluginReactPreamble }}
		<script type="module" src="{{ .ViteURL }}/@vite/client"></script>
		{{- if ne .ViteEntry "" }}
			<script type="module" src="{{ .ViteURL }}/{{ .ViteEntry }}"></script>
		{{- else }}
			<script type="module" src="{{ .ViteURL }}/src/main.tsx"></script>
		{{- end }}
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
