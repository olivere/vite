# Example - Template Registration

This project demonstrates integrating a React application with Go for serving HTML templates. It uses Vite for the frontend build process.

```sh
npm create vite@latest example -- --template react-ts
```

## Configure Vite

We changed the `vite.config.ts` to add the generation of the manifest file and made sure to overwrite the main entry point. Here's how the `vite.config.ts` looks after the changes:

```ts
import react from '@vitejs/plugin-react'
import { defineConfig } from 'vite'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  build: {
    // generates .vite/manifest.json in outDir
    manifest: true,

    rollupOptions: {
      // overwrite default .html entry
      input: "/src/main.tsx",
    },
  },
})
```

## Server side

We then added the [`main.go`](https://github.com/olivere/vite/tree/main/examples/template-registry/main.go).

### Registering Templates

To register a template with the Go server, use the `RegisterTemplate` method of the Handler. Here's an example template and how to register it:

**Example Template**:

```go
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
```

### Register the Template

In the setup code, instantiate the viteHandler with the vite.NewHandler function, passing a configuration. Then you register the customIndex template with the name "index.html".

**Example Set Up**:

```go
  viteHandler, err := vite.NewHandler(vite.Config{
    FS:      os.DirFS("."),
    IsDev:   true,
    ViteURL: "http://localhost:5173",
  })
  if err != nil {
    panic(err)
  }

  viteHandler.RegisterTemplate("index.html", customIndex)
```

### Development mode

If you want to try development mode, first run a new console and do `npm run dev` in the background: It should start the Vite development server on `http://localhost:5173`.

Now run the Go code as:

```sh
$ go run main.go -dev
Listening on on http://127.0.0.1:62002
```

Open up the URL in your browser and you should see the React app, being rendered by a Go HTML template. Not convinced? Open up development mode and go to the `Console`. You should see a message there, which was embedded by the Go code that rendered the HTML.

Notice that you can now change the HTML and JavaScript/TypeScript code, and Hot Module Reload (HMR) should run just fine and update the page inline.

### Production mode

First make sure to run `npm run build` before using production mode, as the Go code relies on embedding the `dist` directory into the Go binary.

Next, simply run the Go code:

```sh
$ go run main.go
Listening on on http://127.0.0.1:61736
```

Open the URL in your browser, and you're seeing a Go template being rendered with an underlying React app.
