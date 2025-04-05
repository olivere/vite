---
prev:
  text: 'Getting Started'
  link: '/guide/getting-started'
next:
  text: 'Examples'
  link: '/guide/examples'
---

# Usage

This page describes the two main approaches for integrating Vite with your Go backend.

[[toc]]

## Option 1: Helper Function

::: info HELPER FUNCTION
The helper function generates the necessary HTML tags (`script`, `link`) to connect your Go application to Vite assets:

- In **development mode**: Links to the Vite dev server
- In **production mode**: Links to your built assets
:::

### Basic Setup

```go
// Initialize the helper function
viteFragment, err := vite.HTMLFragment(vite.Config{
    FS:        os.DirFS("frontend/dist"), // Required: Vite build output directory
    IsDev:     *isDev,                     // Required: Development or Production mode
    ViteURL:   "http://localhost:5173",   // Optional: Defaults to this URL
    ViteEntry: "src/main.js",             // Optional: Depends on your frontend setup
})
if err != nil {
    panic(err)
}

// Create a template
tmpl := template.Must(template.New("index").Parse(indexTemplate))

// Serve the template
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    pageData := map[string]interface{}{
        "Vite": viteFragment,
    }

    if err = tmpl.Execute(w, pageData); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
})

// Template with Vite tags
indexTemplate := `
<head>
    <meta charset="UTF-8" />
    <title>My Go Application</title>
    {{ .Vite.Tags }}
</head>
<body>
    <div id="app"></div>
</body>
`
```

### Serving Assets

With the helper function approach, you need to serve assets yourself:

::: warning Asset Serving Requirements

#### Development Mode

- The Vite dev server handles JS and CSS assets
- Your Go server must serve other assets like images, SVGs, and files in the `public` folder

It's often simpler to [disable](https://vite.dev/config/shared-options.html#publicdir) this folder, as its use case is not primarily for apps with a backend like Go.

#### Production Mode

- You must serve the built assets from the `dist` directory
:::

```go
if isDev {
    // Serve assets from the source directory in development
    serveStaticFolder(mux, "/src/assets/", os.DirFS("frontend/src/assets"))
    // Serve the public folder
    serveStaticFolder(mux, "/", os.DirFS("frontend/public"))
}

if !isDev {
    // Serve assets from the build directory in production
    serveStaticFolder(mux, "/assets/", os.DirFS("frontend/dist/assets"))
    // Serve the public folder
    serveStaticFolder(mux, "/", os.DirFS("frontend/dist"))
}

// Helper function to serve static files
func serveStaticFolder(mux *http.ServeMux, path string, fs fs.FS) {
    mux.Handle(path, http.StripPrefix(path, http.FileServer(http.FS(fs))))
}
```

::: tip FRAMEWORK INTEGRATION
Many Go web frameworks provide simplified methods for serving static files.
For example, with Echo:

```go
e.Static("/assets", "frontend/dist/assets")
```

:::

## Option 2: HTTP Handler

::: info COMPLETE HANDLER
This library provides a complete HTTP handler implementation that can be integrated into your Go server.
:::

### Development Mode

```go
// Create a handler in development mode
handler, err := vite.NewHandler(vite.Config{
    FS:        os.DirFS("./frontend"),        // Source directory of your frontend
    IsDev:     true,                          // Enable development mode
    PublicFS:  os.DirFS("./frontend/public"), // Optional: Serve public directory
    ViteURL:   "http://localhost:5173",       // Optional: Dev server URL
    ViteEntry: "src/main.js"                  // Optional: Entry point
})
if err != nil {
    panic(err)
}

// Use the handler
http.Handle("/", handler)
```

::: warning DEV SERVER REQUIREMENT
In development mode, you still need to run the Vite dev server separately:

```bash
# In a separate terminal
cd frontend
npm run dev
```

:::

### Production Mode

::: tip EMBEDDED ASSETS
In production, you'll typically embed the built Vite assets into your Go binary:
:::

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

// Create a handler in production mode
handler, err := vite.NewHandler(vite.Config{
    FS:    DistFS(),  // Embedded dist directory
    IsDev: false,     // Disable development mode
})
if err != nil {
    panic(err)
}

// Use the handler
http.Handle("/", handler)
```

## Running Your Application

::: details Development Workflow

1. Run the Vite dev server:

   ```bash
   cd frontend
   npm run dev
   ```

2. Run your Go application with development mode enabled:

   ```bash
   go run main.go -dev
   ```

:::

::: details Production Workflow

1. Build your Vite application:

   ```bash
   cd frontend
   npm run build
   ```

2. Run your Go application in production mode:

   ```bash
   go run main.go
   ```

:::

## Configuration Options

A complete list of all configuration parameters for the `vite.Config`.

| Field | Type | Description | Required | Default Value |
|-------|------|-------------|----------|---------------|
| FS | fs.FS | Filesystem containing the Vite assets and manifest. In production, this is the Vite output directory (usually "dist"). In development, this is typically the root directory of the Vite app. | Yes | None |
| PublicFS | fs.FS | Filesystem to serve public files from (usually the "public" directory). Only used in development mode. | No | None |
| IsDev | bool | Determines whether to link to dev Vite server or built assets in 'prod' mode. | Yes | `false` |
| ViteEntry | string | Entrypoint for the Vite application. Useful for implementing secondary routes as described in the Multi-Page App section of the Vite guide. | No | `src/main.tsx` |
| ViteURL | string | Local URL for the Vite development server. Only used in development mode. | No | `http://localhost:5173` |
| ViteManifest | string | File path of the manifest file (relative to FS). Only used in production mode. | No | `.vite/manifest.json` |
| ViteTemplate | Scaffolding | Enum type that specifies which frontend framework is being used. This determines if framework-specific code (preamble) needs to be injected. | No | None |
| AssetsURLPrefix | string | URL prefix for serving asset files. Only used in production mode to construct paths for assets based on the Vite manifest. Useful when serving multiple builds from different base paths. | No | `""` (empty string) |

### ViteTemplate and Preamble

::: info FRAMEWORK INTEGRATION
The `ViteTemplate` option is needed when:

1. You are running in development mode (`IsDev` is `true`)
2. You are using a frontend framework that requires special setup for Hot Module Replacement (HMR)
:::

A **preamble** is a special code snippet injected into the HTML that enables framework-specific features during development. For example, React requires a specific preamble to enable Fast Refresh (Hot Module Replacement). The preamble is a JavaScript snippet that:

- Imports the React Refresh Runtime from the Vite development server
- Injects it into the global window hook
- Sets up necessary refresh registration functions

Without this preamble, React components would not update in real-time during development when you make changes to your code. The library automatically adds the correct preamble based on your `ViteTemplate` setting.

### Scaffolding Options

The `ViteTemplate` field accepts a `Scaffolding` enum with the following values:

| Value | Description | Requires Preamble |
|-------|-------------|-------------------|
| React | Template for a React project | ✅ |
| ReactTs | Template for a TypeScript React project | ✅ |
| ReactSwc | Template for a React project using SWC compiler | ✅ |
| ReactSwcTs | Template for a TypeScript React project using SWC | ✅ |
| Vue | Template for a Vue.js project | ❌ |
| VueTs | Template for a TypeScript Vue.js project | ❌ |
| Vanilla | Template for a Vanilla JavaScript project | ❌ |
| VanillaTs | Template for a Vanilla TypeScript project | ❌ |
| Preact | Template for a Preact project | ❌ |
| PreactTs | Template for a TypeScript Preact project | ❌ |
| Lit | Template for a Lit project | ❌ |
| LitTs | Template for a TypeScript Lit project | ❌ |
| Svelte | Template for a Svelte project | ❌ |
| SvelteTs | Template for a TypeScript Svelte project | ❌ |
| Solid | Template for a Solid project | ❌ |
| SolidTs | Template for a TypeScript Solid project | ❌ |
| Qwik | Template for a Qwik project | ❌ |
| QwikTs | Template for a TypeScript Qwik project | ❌ |
| None | Opt out of using a specific scaffolding | ❌ |

### Additional Notes

::: tip

- React-based templates (React, ReactTs, ReactSwc, ReactSwcTs) require a preamble for Hot Module Replacement (HMR) to work properly in development mode.
- The `PublicFS` field is optional. If not provided, the system will check if the "public" directory exists in the Vite app and serve files from there.

:::

::: warning

- In development mode, the `ViteURL` parameter defines the base URL for assets, making the `AssetsURLPrefix` parameter unnecessary.
- The manifest file is used in production mode to map original file paths to transformed file paths.

:::