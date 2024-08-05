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
	customIndex = `
<!doctype html>
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
    <header>My Custom Header</header>
    <main>
      <div id="root"></div>
      <noscript>
        <div>
          <h2>JavaScript Required</h2>
          <p>This application requires JavaScript to run. Please enable JavaScript
          in your browser settings and reload the page.</p>
        </div>
      </noscript>
    </main>
    <footer>My Unique Footer</footer>
  </body>
</html>
	`
)

func runDevServer() {
	// Handle the Vite server.
	viteHandler, err := vite.NewHandler(vite.Config{
		FS:      os.DirFS("."),
		IsDev:   true,
		ViteURL: "http://localhost:5173",
	})
	if err != nil {
		panic(err)
	}

	viteHandler.RegisterTemplate("index.html", customIndex)

	// Create a new handler.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		Handler: handler,
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

	// Create a new handler.
	viteHandler, err := vite.NewHandler(vite.Config{
		FS:    fs,
		IsDev: false,
	})
	if err != nil {
		panic(err)
	}

	viteHandler.RegisterTemplate("index.html", customIndex)

	// Create a new handler.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		Handler: handler,
	}

	log.Printf("Listening on on http://%s", l.Addr())

	// Start the server.
	if err := server.Serve(l); err != nil {
		panic(err)
	}
}
