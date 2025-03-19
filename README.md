# Vite Backend integration with Go

This library facilitates the integration of a Vite-based frontend with a Go-based backend, following the guidelines outlined in the [Vite backend integration guide](https://vitejs.dev/guide/backend-integration.html).

> [!NOTE]
> Please follow the [Vite guidelines](https://vitejs.dev/guide/backend-integration.html) to configure your Vite-based frontend, i.e. `vite.config.(js|ts)`. E.g. you need to make sure that the `manifest.json` is being generated for production.

## Usage as a simple 'helper function'

In this setup, this library can be used as a helper function which generates HTML tags (`script`, `link`) to point to your Vite assets. In development mode (when the Vite server is running) this links to the Vite server, and in production mode, this links to your built assets. You are responsible for serving those built assets in production (and some assets like images in dev mode) however it is simple to do and many golang web frameworks have methods to make it even easier ([see asset serving](#serving-assets)).

We setup the helper function, and then pass it to your HTML template.

```go
viteFragment, err := vite.HTMLFragment(vite.Config{
    FS:      os.DirFS("frontend/dist"), // required: Vite build output
    IsDev:   *isDev,                    // required: true or false
    ViteURL: "http://localhost:5173",   // optional: defaults to this
    ViteEntry: "src/main.js",           // optional: dependent on your frontend stack
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

This example should give you an idea of how to use this in your application. It is designed to be as simple as possible and independent of your framework, you just need to specify some config and then call `viteFragment.Tags` in your template. See the list of [examples](#examples) to get started.

### Serving Assets

The code above only produces the HTML tags. You are responsible for serving assets as this varies depending on your framework and setup. For example, you may or may not want to use the `public` folder in Vite. If you do use it, you need to serve its contents in dev and prod modes.

You need to setup everything depending on the mode.

#### In development

The Vite web server (normally running at `http://localhost:5173`) serves `js` and `css` assets for you so you do **not** need to worry about serving those in your Go backend. However, if you have assets like `svg`, images etc. that you are referencing in your frontend, you need to serve these in development yourself. The classical example of this is the `public` folder. You are also responsible for serving the `public` folder in both dev and prod.

> [!NOTE]
> It's often simpler to [disable](https://vite.dev/config/shared-options.html#publicdir) this folder, as its use case is not primarily for apps with a backend like Go.)

#### In production

A Vite build (by default) produces output in the `dist` folder. You will want to serve the `dist/assets` directory wherever it is located in your project under the `/assets/` URL. With a default Vite config, your script tags will be looking for an entry point like `/assets/main-BLC8vTVb.js`. This library produces its HTML tags based on the `dist/.vite/manifest.json`, and this manifest file drives the setup.

An example of serving the assets using `net/http` would be:

```go
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

#### Development mode

In development mode, you need to run Vite with something like `npm run dev`, `pnpm dev` etc.

It is probably best to setup a flag in your go app to swap between dev and prod modes, which is easy with `flag` from the standard library, or use an environment variable like `ENV` for that.

```go
var isDev = flag.Bool("dev", false, "run in development mode")
flag.Parse()
```

Then in a separate shell from Vite, you can pass that flag to start the app.

```sh
go run main.go -dev
```

#### Production mode

In production mode, it's even simpler. You run a Vite build to generate the assets, and the Go binary would embed and serve these:

```sh
npm run build
go run main.go
```

## Usage with the provided Handler

This integration is done by a HTTP handler, implementing `http.Handler`. The handler, again, has two modes: Development and production.

### Development

In development mode, you need to create the handler by passing a file system that points to a source of your Vite app as the first parameter. The second parameter needs to be true to put the handler into development mode. The rest of the parameters let you specify how to integrate the Vite server and entry point as well as the public directory (if any). Notice that in development mode, the Vite server is running in the background, typically `http://localhost:5173` (the endpoint served by running `npm run dev`, `pnpm dev` etc.). If that server is on a different endpoint, you need to configure that with the `ViteURL`. Again: You need to run the Vite server in the background in development mode, so open up a 2nd console and run something like `npm run dev`.

Here's an example of initializing the HTTP handler:

```go
// Serve in development mode (assuming your frontend code is in ./frontend,
// relative to your binary)
v, err := vite.NewHandler(vite.Config{
    FS:        os.DirFS("./frontend"),
    IsDev:     true,
    PublicFS:  os.DirFS("./frontend/public"), // optional: we use the "public" directory under "FS" by default
    ViteURL:   "http://localhost:5173",       // optional: we use "http://localhost:5173" by default
    ViteEntry: "src/main.js"                  // optional: depending on your frontend stack
})
if err != nil { ... }
```

### Production

In production mode, you typically embed the whole generated `dist` directory generated by `vite build` into the Go binary, using `go:embed`. In that case, your `FS` parameter needs to be the embedded `dist` file system. The `IsDev` parameter must be `false` to enable production mode. The other parameters are optional in production mode.

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

## Configuration

Here's a complete list of all configuration parameters of the `vite.Config`.

| Field        | Type                                                                            | Description                                                                                                                                                             | Useful Default                  |
|--------------|---------------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------|---------------------------------|
| IsDev        | bool                                                                            | Instruct whether to link to dev Vite server or built assets in 'prod'                                                                                                   | `false`                         |
| FS           | fs.FS                                                                           | FS containing the Vite assets (and manifest)                                                                                                                            |                                 |
| ViteEntry    | string                                                                          | (optional) Entrypoint for the Vite application. Usually a main Javascript file. This is the top of the dependency tree and Vite will import dependencies based on this entrypoint. | `src/main.tsx`                  |
| ViteURL      | string                                                                          | (optional) Local URL for the Vite development server. Not used in production mode.                                                                                                 | `http://localhost:5173`         |
| ViteManifest | string                                                                          | (optional) File path of the manifest file (relative to FS). Only used in production mode.                                                                                          | `.vite/manifest.json`           |
| ViteTemplate | [Scaffolding](https://github.com/olivere/vite/blob/main/config.go#L53C6-L53C17) | (optional) A enum-like type that instruct this library what preambles to inject based on what project type (React, Vue etc). Needed for React applications to enable HMR etc.      | React (includes React preamble) |
| PublicFS     | fs.FS                                                                           | (optional) Only used with the `vite.NewHandler` to serve the public directory for you. Not used when this library is used as a helper template function.                       |                                 |

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

### Router App

This application consists of a Go backend, serving a Vite-based app using TanStack Router and TanStack Query libraries. See the the [`examples/router` directory](https://github.com/olivere/vite/tree/main/examples/router).

## License

See license in LICENSE file.
