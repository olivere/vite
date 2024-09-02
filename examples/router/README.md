# Example

This application is created as [described here](https://vitejs.dev/guide/):

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

We then added the [`main.go`](https://github.com/olivere/vite/tree/main/examples/basic/main.go).

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
