package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/olivere/vite"
)

//go:embed all:dist
var dist embed.FS

//go:embed index.tmpl
var goIndex embed.FS

func main() {
	var isDev = flag.Bool("dev", false, "run in development mode")

	flag.Parse()

	// The following block sets up the environment and configuration for a Go
	// application integrated with Vite. It determines whether the application
	// should run in development or production mode based on the 'isDev' flag

	var viteAssetsDir string
	var viteAssetsURL string
	var viteFS fs.FS

	if *isDev {
		viteAssetsDir = "src/assets"
		viteAssetsURL = "/src/assets/"
		viteFS = os.DirFS(".")
	} else {
		viteAssetsDir = "dist/assets"
		viteAssetsURL = "/assets/"
		fs, err := fs.Sub(dist, "dist")
		if err != nil {
			panic(err)
		}
		viteFS = fs
	}

	mux := http.NewServeMux()

	// Serve assets that Vite would treat as 'public' assets.
	//
	// In this example, the 'static' directory is used as a replacement for
	// Vite's default 'public' folder, and 'publicDir' is disabled in the Vite
	// config. We're using the 'static' directory to achieve similar
	// functionality, but available to the Go backend and Vite.
	//
	// To use a static asset in our Vite frontend, we import it like this:
	//
	// import viteLogo from '/static/vite.svg'

	staticAssets := http.FileServer(http.Dir("static"))

	mux.Handle("/static/", http.StripPrefix("/static/", staticAssets))

	// Serve Vite-managed assets from the Go backend, accommodating both
	// development and production environments.
	//
	// Usage in Vite remains the same as in a standard Vite setup. The Go backend
	// will serve the assets from the correct location based on the environment.

	viteAssets := http.FileServer(http.Dir(viteAssetsDir))

	mux.Handle(viteAssetsURL, http.StripPrefix(viteAssetsURL, viteAssets))

	// This block demonstrates the setup and usage of `vite.Middleware` in a Go
	// web application.
	//
	// The Middleware is then applied to the root route ("/") of an HTTP mux,
	// where it processes incoming requests before rendering a "index.tmpl"
	// template.
	//
	// This setup allows a minimal integration of Vite with the Go web server,
	// handling both development and production environments.

	viteMiddleware, err := vite.NewMiddleware(vite.Config{
		FS:      viteFS,
		IsDev:   *isDev,
		ViteURL: "http://localhost:5173",
	})
	if err != nil {
		panic(err)
	}

	mux.HandleFunc("/", viteMiddleware.Use(func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index.tmpl").ParseFS(goIndex, "index.tmpl")

		if err != nil {
			panic(err)
		}

		if err = tmpl.Execute(w, nil); err != nil {
			panic(err)
		}
	}))

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
