# Vite Backend integration with Go

> **Note**
> This library is still a work in progress and may not be stable or fully functional. Use it at your own risk.

This library implements a [Vite backend integration](https://vitejs.dev/guide/backend-integration.html) for Go. Please follow the guidelines there to configure your Vite project, i.e. `vite.config.(js|ts)`. E.g. you need to make sure that the `manifest.json` is being generated for production.

## Usage as a simple 'helper function'

In this setup, this library can be used as a helper function which generates HTML tags (`script`, `link`) to point to your Vite assets. In development mode (when the Vite server is running) this links to the Vite server, and in production mode, this links to your built assets. You are responsible for serving those built assets in production (and some assets like images in dev mode) however it is simple to do and many golang web frameworks have methods to make it even easier. [See asset serving](#serving-assets)

We setup the helper function, and then pass it to your HTML template.

```golang
viteFragment, err := vite.HTMLFragment(vite.Config{
    FS:      os.DirFS("frontend/dist"),
    IsDev:   *isDev,
    ViteURL: "http://localhost:5173",
})
if err != nil {
    panic(err)
}

t := template.Must(template.New("name").Parse(indexTemplate))

http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    pageData := map[string]interface{}{
        "Vite": viteFragment,
    }

    if err = tmpl.Execute(w, pageData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
  
indexTemplate := `
<head>
    <meta charset="UTF-8" />
    <title>My Go Application</title>
    {{ .Vite.Tags }}
</head>
<body></body>
`
```

This example should give you an idea of how to use this in your application. It is designed to be as simple as possible and independent of your framework, you just need to specify some config and then call `viteFragment.Tags` in your template. For a full working example, see [examples](#examples)

### Serving Assets
The code above only produces the HTML tags. You are responsible for serving assets as this can vary hugely depending on your framework and setup. For example, you may want to use the `public` folder in Vite, or it is often simpler to disable it. If you do use it, you need to serve its contents in dev and prod modes.

You need to setup:

**In production:**
- Vite produces a `dist` folder, you will want to serve `dist/assets` wherever it is located in your project under the `/assets/` URL. With a default Vite config, your script tags will be looking for a resource like `/assets/main-BLC8vTVb.js`. This library produces its HTML tags based on the `dist/.vite/manifest.json` so this drives the setup.

**In development**
- The Vite web server (normally running at `http://localhost:5173`) serves `js` and `css` assets for you so you do **not** need to worry about serving those. However, if you have assets like `svg`, images etc **that you are importing** into your Javascript app, you need to serve these in development yourself. You are also responsible for serving the `public` folder in dev and prod. (It's simpler to [disable](https://v2.vitejs.dev/config/#publicdir) this folder, its use case is not primarily for apps with a backend like Go). 

A simple example of serving the assets using `net/http` would be:
```golang
if *isDev {
    serverStaticFolder(mux, "/src/assets/", os.DirFS("frontend/src/assets"))
} else {
    serverStaticFolder(mux, "/assets/", os.DirFS("frontend/dist/assets"))
}

func serverStaticFolder(mux *http.ServeMux, path string, fs fs.FS) {
	mux.Handle(path, http.StripPrefix(path, http.FileServerFS(fs)))
}
```

If you are using a framework like [Echo](https://echo.labstack.com/) then it provides [this functionality already](https://echo.labstack.com/docs/static-files#using-echostatic).

### Running your app
For development mode:

You need to run Vite with something like `npm run dev`, `pnpm dev` etc.

It is probably best to setup a flag in your go app to swap between dev and prod modes, which is easy with `flag` from the standard library.
```golang
var isDev = flag.Bool("dev", false, "run in development mode")
flag.Parse()
```
Then in a separate shell from Vite, you can run `go run main.go -dev` to start the app.

In production *mode* it's easier:
```
$ npm run build
$ go run main.go
```
## Usage with the provided Handler

This integration is done by a HTTP handler, implementing `http.Handler`. The handler has two modes: Development and production.

In development mode, you need to create the handler by passing a file system that points to a source of your Vite app as the first parameter. The second parameter needs to be true to put the handler into development mode. And the third and last parameter points to the Vite server running in the background, typically `http://localhost:5173` (the endpoint served by running `npm run dev`, `pnpm dev` etc.). Again: You need to run the Vite server in the background in development mode, so open up a 2nd console and run something like `npm run dev`.

Example:

```go
// Serve in development mode (assuming your frontend code is in ./frontend,
// relative to your binary)
v, err := vite.NewHandler(vite.Config{
    FS:       os.DirFS("./frontend"),
    IsDev:    true,
    PublicFS: os.DirFS("./frontend/public"), // optional: we use the "public" directory under "FS" by default
    ViteURL:  "http://localhost:5173",       // optional: we use "http://localhost:5173" by default
})
if err != nil { ... }
```

In production mode, you typically embed the whole generated dist directory generated by `vite build` into the Go binary, using `go:embed`. In that case, your first parameter needs to be the embedded "dist" file system. The second parameter must be false to enable production mode. The last parameter can be blank, as it is not used in production mode.

Example:

```go
//go:embed all:dist
var distFS embed.FS

func DistFS() fs.FS {
    efs, err := fs.Sub(distFS, "dist")
    if err != nil {
        panic(fmt.Sprintf("unable to serve frontend: %v", err))
    }
    return efs
}

// Serve in production mode
v, err := vite.NewHandler(vite.Config{
    FS:    DistFS(),
    IsDev: false,
})
if err != nil { ... }
```

## Examples

### Simple Helper Function 

See the [`examples/helper-function-basic` directory](https://github.com/olivere/vite/tree/main/examples/helper-function-basic) for a demonstration of a very basic React app that integrates a Go backend.

### Basic with Handler

See the [`examples/basic` directory](https://github.com/olivere/vite/tree/main/examples/basic) for a demonstration of a very basic React app that integrates a Go backend.

### Inertia.js

[See this example](https://github.com/danclaytondev/go-inertia-vite) for using Golang with `net/http`, [Inertia.js](https://inertiajs.com/) and this library for managing Vite assets.

### Multi Page App

For Vite apps that have multiple entry points, you can pass the entry point by creating a separate `vite.Handler` and specifying the `ViteEntry` field. See the [`examples/multi-page-app` directory](https://github.com/olivere/vite/tree/main/examples/multi-page-app) for an example.

### Template Registration

You can use custom HTML templates in your Go backend for serving different React pages. See the [`examples/template-registry` directory](https://github.com/olivere/vite/tree/main/examples/template-registry) for an example.

## License

See license in LICENSE file.
