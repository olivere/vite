package vite

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
)

// Fragment holds HTML content generated for Vite integration, intended to be
// embedded in HTML templates.
type Fragment struct {
	// Tags contains a string of HTML content for embedding Vite-related assets
	// such as JavaScript and CSS. The content is stored as template.HTML to
	// ensure it is rendered without escaping within the HTML template.
	Tags template.HTML
}

// HTMLFragment generates an HTML fragment for Vite integration based on the provided configuration.
//
// This function takes a Config struct and uses it to create the necessary HTML
// elements for including Vite-managed assets in a web page. The returned HTML
// is typically intended to be placed in the <head> section of an HTML document.
//
// Parameters:
//   - config: A Config struct containing the necessary configuration options for Vite integration.
//
// Returns:
//   - template.HTML: An HTML fragment as a template.HTML type, which can be safely embedded in HTML templates.
//   - error: An error if the HTML fragment generation fails for any reason.
//
// The returned HTML fragment may include elements such as:
//   - Script tags for Vite client in development mode
//   - Links to stylesheets
//   - Module preload tags
//   - Script tags for entry points and other necessary JavaScript
//
// Usage example:
//
//	fragment, err := HTMLFragment(myConfig)
//	if err != nil {
//	    // Handle error
//	}
//	// Use fragment in your HTML template
func HTMLFragment(config Config) (*Fragment, error) {
	pd := &pageData{
		IsDev:     config.IsDev,
		ViteEntry: config.ViteEntry,
		ViteURL:   config.ViteURL,
	}

	if config.IsDev {
		// Development mode.
		if pd.ViteURL == "" {
			pd.ViteURL = "http://localhost:5173"
		}

		// Check if the specified Vite template requires a preamble and set the
		// corresponding preamble string in the plugin configuration.
		//
		// If the Vite template value is less than 1, it is considered as an
		// uninitialized state, and the default React preamble is applied.
		// Otherwise, if the template requires a preamble, it uses the
		// specific preamble for the given Vite template.
		if config.ViteTemplate < 1 {
			pd.PluginReactPreamble = template.HTML(React.Preamble(pd.ViteURL))
		} else if config.ViteTemplate.RequiresPreamble() {
			pd.PluginReactPreamble = template.HTML(config.ViteTemplate.Preamble(pd.ViteURL))
		}
	} else {
		if config.ViteManifest == "" {
			config.ViteManifest = ".vite/manifest.json"
		}
		mf, err := config.FS.Open(config.ViteManifest)
		if err != nil {
			return nil, fmt.Errorf("vite: open manifest: %w", err)
		}
		defer mf.Close()

		m, err := ParseManifest(mf)
		if err != nil {
			return nil, fmt.Errorf("vite: parse manifest: %w", err)
		}
		var chunk *Chunk
		if pd.ViteEntry == "" {
			chunk = m.GetEntryPoint()
		} else {
			entries := m.GetEntryPoints()
			for _, entry := range entries {
				if pd.ViteEntry == entry.Src {
					chunk = entry
					break
				}
			}
		}
		if chunk == nil {
			return nil, fmt.Errorf("vite: unable to find chunk for entry point %q", pd.ViteEntry)
		}

		pd.StyleSheets = template.HTML(m.GenerateCSS(chunk.Src, config.AssetsURLPrefix))
		pd.Modules = template.HTML(m.GenerateModules(chunk.Src, config.AssetsURLPrefix))
		pd.PreloadModules = template.HTML(m.GeneratePreloadModules(chunk.Src, config.AssetsURLPrefix))
	}

	// Create a buffer to store the executed template output
	var buf bytes.Buffer

	// Pass the JoinPath function to the template so we
	// can use {{ urljoin .base .path }}
	templateFuncs := template.FuncMap{
		"urljoin": url.JoinPath,
	}

	// Parse the predefined headTmpl into a new template
	tmpl, err := template.New("vite").Funcs(templateFuncs).Parse(htmlTmpl)
	if err != nil {
		// Return an error if parsing fails
		return nil, fmt.Errorf("vite: parse template: %w", err)
	}

	// Execute the template with pd (PageData) as the data source
	err = tmpl.Execute(&buf, pd)
	if err != nil {
		// Return an error if template execution fails
		return nil, fmt.Errorf("vite: execute template: %w", err)
	}

	return &Fragment{Tags: template.HTML(buf.Bytes())}, nil
}

// htmlTmpl is a constant string that contains a Go template for including
// Vite-related scripts and stylesheets in a <head> element of an HTML page.
// This template adapts its output based on whether the application is running
// in development or production mode.
const htmlTmpl = `
{{- if .IsDev }}
	{{ .PluginReactPreamble }}
	<script type="module" src="{{ urljoin .ViteURL "/@vite/client" }}"></script>
	{{- if ne .ViteEntry "" }}
		<script type="module" src="{{ urljoin .ViteURL .ViteEntry }}"></script>
	{{- else }}
		<script type="module" src="{{ urljoin .ViteURL "/src/main.tsx" }}"></script>
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
