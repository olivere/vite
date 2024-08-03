package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/olivere/vite"
)

//go:embed all:dist
var dist embed.FS

func main() {
	var (
		isDev = flag.Bool("dev", false, "run in development mode")
	)
	flag.Parse()

	if *isDev {
		runDevServer()
	} else {
		runProdServer()
	}
}

var (
	nestedHTML = `<!doctype html>
<html lang="en" class="h-full scroll-smooth">
  <head>
    <meta charset="UTF-8" />
	{{- if .Metadata }}
		{{ .Metadata }}
	{{- end }}
	{{- if .IsDev }}
		{{ .PluginReactPreamble }}
		<script type="module" src="{{ .ViteURL }}/@vite/client"></script>
		{{ if ne .ViteEntry "" }}
			<script type="module" src="{{ .ViteURL }}/{{ .ViteEntry }}"></script>
		{{ else }}
			<script type="module" src="{{ .ViteURL }}/src/main.tsx"></script>
		{{ end }}
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

func runDevServer() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle the Vite server.
		viteHandler, err := vite.NewHandler(vite.Config{
			FS:      os.DirFS("."),
			IsDev:   true,
			ViteURL: "http://localhost:5173",
		})
		if err != nil {
			panic(err)
		}

		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			// Server the index.html file.
			ctx := r.Context()
			ctx = vite.MetadataToContext(ctx, vite.Metadata{
				Title: "Hello, Vite!",
			})
			ctx = vite.ScriptsToContext(ctx, `<script>console.log('Hello, nice to meet you in the console!')</script>`)
			viteHandler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		viteHandler.ServeHTTP(w, r)
	})

	mux.HandleFunc("/nested", func(w http.ResponseWriter, r *http.Request) {
		// Handle the Vite server.
		viteHandler, err := vite.NewHandler(vite.Config{
			FS:        os.DirFS("."),
			IsDev:     true,
			ViteEntry: "src/nested.tsx",
			ViteURL:   "http://localhost:5173",
		})
		if err != nil {
			panic(err)
		}

		viteHandler.RegisterTemplate("/nested", nestedHTML)

		if r.URL.Path == "/nested" {
			// Server the index.html file.
			ctx := r.Context()
			ctx = vite.MetadataToContext(ctx, vite.Metadata{
				Title: "Hello, Nested Vite!",
			})
			ctx = vite.ScriptsToContext(ctx, `<script>console.log('Hello Nested!, nice to meet you in the console!')</script>`)
			viteHandler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		viteHandler.ServeHTTP(w, r)
	})

	// Start a listener.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		var err1 error
		if l, err1 = net.Listen("tcp6", "[::1]:0"); err1 != nil {
			panic(fmt.Errorf("starting HTTP server: %v", err))
		}
	}

	// Create a new server.
	server := http.Server{
		Handler: mux,
	}

	log.Printf("Listening on on http://%s", l.Addr())

	// Start the server.
	if err := server.Serve(l); err != nil {
		panic(err)
	}
}

func runProdServer() {
	fs, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Create a new handler.
		viteHandler, err := vite.NewHandler(vite.Config{
			FS:    fs,
			IsDev: false,
		})
		if err != nil {
			panic(err)
		}

		if r.URL.Path == "/" || r.URL.Path == "/index.html" {
			// Server the index.html file.
			ctx := r.Context()
			ctx = vite.MetadataToContext(ctx, vite.Metadata{
				Title: "Hello, Vite (Prod)!",
			})
			ctx = vite.ScriptsToContext(ctx, `<script>console.log('Hello, nice to meet you in the console!')</script>`)
			viteHandler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		viteHandler.ServeHTTP(w, r)
	})

	mux.HandleFunc("/nested", func(w http.ResponseWriter, r *http.Request) {

		// Create a new handler.
		viteHandler, err := vite.NewHandler(vite.Config{
			FS:        fs,
			IsDev:     false,
			ViteEntry: "src/nested.tsx",
		})
		if err != nil {
			panic(err)
		}

		viteHandler.RegisterTemplate("/nested", nestedHTML)

		if r.URL.Path == "/nested" {
			// Server the index.html file.
			ctx := r.Context()
			ctx = vite.MetadataToContext(ctx, vite.Metadata{
				Title: "Hello, Nested Vite (Prod)!",
			})
			ctx = vite.ScriptsToContext(ctx, `<script>console.log('Hello Nested, nice to meet you in the console!')</script>`)
			viteHandler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		viteHandler.ServeHTTP(w, r)
	})

	// Start a listener.
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		var err1 error
		if l, err1 = net.Listen("tcp6", "[::1]:0"); err1 != nil {
			panic(fmt.Errorf("starting HTTP server: %v", err))
		}
	}

	// Create a new server.
	server := http.Server{
		Handler: mux,
	}

	log.Printf("Listening on on http://%s", l.Addr())

	// Start the server.
	if err := server.Serve(l); err != nil {
		panic(err)
	}
}
