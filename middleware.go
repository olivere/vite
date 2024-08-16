package vite

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"log/slog"
	"net/http"
)

// Middleware integrates Vite with a Go application, offering functionality
// to serve assets processed by Vite and manage HTTP requests in both
// development and production environments.
//
// It manages configuration, asset manifests, and serves static files
// using various file system interfaces.
type Middleware struct {
	config     *Config
	manifest   *Manifest
	pub        fs.FS
	pubFS      http.FileSystem
	pubHandler http.Handler
}

// NewMiddleware initializes a new Middleware instance with the specified
// configuration. It sets up the necessary components to integrate Vite
// with a Go application for serving assets and handling HTTP requests.
//
// Parameters:
//   - config: The configuration to use for setting up the Middleware.
//
// Returns:
//   - A pointer to the newly created Middleware instance and an error, if any
//     occurred during initialization.
func NewMiddleware(config Config) (*Middleware, error) {
	if config.FS == nil {
		return nil, fmt.Errorf("vite: fs is nil")
	}

	m := &Middleware{
		config: &config,
	}

	if m.config.IsDev == false {
		mf, err := config.FS.Open(".vite/manifest.json")
		if err != nil {
			return nil, fmt.Errorf("vite: open manifest: %w", err)
		}
		defer mf.Close()

		m.manifest, err = ParseManifest(mf)
		if err != nil {
			return nil, fmt.Errorf("vite: parse manifest: %w", err)
		}
	} else {
		if config.ViteURL == "" {
			m.config.ViteURL = "http://localhost:5173"
		}

		if config.PublicFS == nil {
			pub, err := fs.Sub(config.FS, "public")
			if err != nil {
				return nil, fmt.Errorf("vite: default public folder: %w", err)
			}
			m.pub = pub
			m.pubFS = http.FS(m.pub)
			m.pubHandler = http.FileServerFS(m.pub)
		} else {
			m.pub = config.PublicFS
			m.pubFS = http.FS(m.pub)
			m.pubHandler = http.FileServerFS(m.pub)
		}
	}

	return m, nil
}

type customResponseWriter struct {
	http.ResponseWriter
	body []byte
}

func (crw *customResponseWriter) Write(b []byte) (int, error) {
	crw.body = append(crw.body, b...)
	return len(b), nil
}

const viteTmpl = `
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
`

// insertViteHTML inserts Vite-related HTML into a byte slice at a specified marker.
//
// This function is used to inject Vite's development scripts or production
// asset links into an HTML template.
// /
// Parameters:
//   - content: The original HTML content as a byte slice.
//   - marker: A string that serves as an insertion point in the content.
//   - html: The HTML string to be inserted before the marker./
//
// Returns:
//   - A new byte slice with the HTML inserted before the marker.
//   - An error if the marker is not found in the content.
//
// The function searches for the first occurrence of the marker in the content.
// If found, it inserts the provided HTML immediately before the marker.
// If the marker is not found, it returns an error.
//
// Note: This function replaces only the first occurrence of the marker.
func insertViteHTML(content []byte, marker, html string) ([]byte, error) {
	mb := []byte(marker)
	if bytes.Index(content, mb) < 0 {
		return nil, fmt.Errorf("vite: template marker not found: %q", marker)
	}
	return bytes.Replace(content, mb, []byte(html+marker), 1), nil
}

// Use wraps the provided HTTP handler function with the Middleware's logic,
// allowing it to process requests before passing them to the next handler.
//
// Parameters:
//   - next: The next http.HandlerFunc in the chain to be called after the
//     Middleware has processed the request.
//
// Returns:
//   - A new http.HandlerFunc that incorporates the Middleware's functionality
//     before delegating to the provided next handler.
func (m *Middleware) Use(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		crw := &customResponseWriter{ResponseWriter: w}

		// Invoke `next` early to generate the parent's response, writing to `crw`.
		next.ServeHTTP(crw, r)

		viteData := pageData{
			IsDev:     m.config.IsDev,
			ViteEntry: m.config.ViteEntry,
			ViteURL:   m.config.ViteURL,
		}

		if m.config.IsDev {
			viteData.PluginReactPreamble = template.HTML(PluginReactPreamble(m.config.ViteURL))
		} else {
			var chunk *Chunk
			if chunk == nil {
				if viteData.ViteEntry == "" {
					chunk = m.manifest.GetEntryPoint()
				} else {
					entries := m.manifest.GetEntryPoints()
					for _, entry := range entries {
						if viteData.ViteEntry == entry.Src {
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
			viteData.StyleSheets = template.HTML(m.manifest.GenerateCSS(chunk.Src))
			viteData.Modules = template.HTML(m.manifest.GenerateModules(chunk.Src))
			viteData.PreloadModules = template.HTML(m.manifest.GeneratePreloadModules(chunk.Src))
		}

		tmpl, err := template.New("vite").Parse(viteTmpl)
		if err != nil {
			slog.Error("vite: parse middleware template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		// Use a buffer to execute the `tmpl`, applying `viteData` to the
		// template's relevant placeholders
		var buf bytes.Buffer

		err = tmpl.Execute(&buf, viteData)
		if err != nil {
			slog.Error("vite: execute middleware template", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}

		resp, err := insertViteHTML(crw.body, "</head>", buf.String())
		if err != nil {
			slog.Error("vite: inserting vite html", err)
			http.Error(w, "Iternal server error", http.StatusInternalServerError)
		}

		w.Write(resp)
	}
}
