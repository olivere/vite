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

//go:embed all:static
var static embed.FS

//go:embed index.tmpl
var goIndex embed.FS

func main() {
	var (
		isDev = flag.Bool("dev", false, "run in development mode")
	)
	flag.Parse()

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

	var staticFileServer http.Handler
	if *isDev {
		staticFileServer = http.FileServer(http.Dir("static"))
	} else {
		staticFS, err := fs.Sub(static, "static")
		if err != nil {
			panic(err)
		}
		staticFileServer = http.FileServer(http.FS(staticFS))
	}

	mux.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// Serve Vite-managed assets from the Go backend, accommodating both
	// development and production environments.
	//
	// Usage in Vite remains the same as in a standard Vite setup. The Go backend
	// will serve the assets from the correct location based on the environment.

	var viteAssetsFileServer http.Handler
	var viteAssetsSrcAttribute string // vite changes this path, depending on dev or prod
	if *isDev {
		viteAssetsSrcAttribute = "/src/assets/"
		viteAssetsFileServer = http.FileServer(http.Dir("src/assets"))
	} else {
		viteAssetsFS, err := fs.Sub(static, "static/dist/assets")
		if err != nil {
			panic(err)
		}
		viteAssetsFileServer = http.FileServer(http.FS(viteAssetsFS))
		viteAssetsSrcAttribute = "/assets/"
	}

	mux.Handle(viteAssetsSrcAttribute, http.StripPrefix(viteAssetsSrcAttribute, viteAssetsFileServer))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// viteFS represents the filesystem used for serving Vite assets.
		//
		// This variable is used to configure the Vite integration, ensuring that
		// the correct assets are served based on the current environment.
		var viteFS fs.FS
		if *isDev {
			viteFS = os.DirFS(".")
		} else {
			fs, err := fs.Sub(static, "static/dist")
			if err != nil {
				panic(err)
			}
			viteFS = fs
		}

		// viteFragment generates the necessary HTML for Vite integration.
		//
		// It calls vite.HTMLFragment with a Config struct to create an HTML fragment
		// that includes all required Vite assets and scripts. This fragment is intended
		// to be inserted into the <head> section of the HTML document.
		viteFragment, err := vite.HTMLFragment(vite.Config{
			FS:      viteFS,
			IsDev:   *isDev,
			ViteURL: "http://localhost:5173",
		})
		if err != nil {
			panic(err)
		}

		tmpl, err := template.New("index.tmpl").ParseFS(goIndex, "index.tmpl")

		if err != nil {
			panic(err)
		}

		if err = tmpl.Execute(w, map[string]interface{}{
			"Vite": viteFragment,
		}); err != nil {
			panic(err)
		}
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

	log.Printf("Listening on http://%s", l.Addr())

	// Start the server.
	if err := server.Serve(l); err != nil {
		panic(err)
	}
}
