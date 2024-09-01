package vite

import "io/fs"

// Config is the configuration for the handler.
type Config struct {
	// FS is the file system to serve files from. In production, this is
	// the Vite output directory, which usually is the "dist" directory.
	// In development, this is usually the root directory of the Vite app.
	FS fs.FS

	// PublicFS is the file system to serve public files from. This is
	// usually the "public" directory. It is optional and can be nil.
	// If it is nil, we will check if the "public" directory exists in
	// the Vite app, and serve files from there. If it does not exist,
	// we will not serve any public files. It is only used in development
	// mode.
	PublicFS fs.FS

	// IsDev is true if the server is running in development mode, false
	// otherwise.
	IsDev bool

	// ViteEntry specifies the path to a particular entry point in the Vite
	// manifest. This is useful for implementing secondary routes, similar to the
	// example provided in the [Multi-Page App] section of the Vite guide.
	//
	// [Multi-Page App]: https://vitejs.dev/guide/build.html#multi-page-app
	ViteEntry string

	// ViteURL is the URL of the Vite server, used to load the Vite client
	// in development mode (and defaults to http://localhost:5173).
	// It is unused in production mode.
	ViteURL string

	// ViteManifest is the path to the Vite manifest file. This is used in
	// production mode to load the manifest file and map the original file
	// paths to the transformed file paths. If this is not provided, the
	// default path is ".vite/manifest.json".
	ViteManifest string

	// ViteTemplate specifies a configuration template used to scaffold the Vite
	// project. See [Scaffolding Your First Vite Project].
	//
	// [Scaffolding Your First Vite Project]: https://vitejs.dev/guide/#scaffolding-your-first-vite-project
	ViteTemplate Scaffolding
}

// Scaffolding represents various templates provided by Vite that can be used
// to scaffold a Vite project. See [Scaffolding Your First Vite Project].
//
// [Scaffolding Your First Vite Project]: https://vitejs.dev/guide/#scaffolding-your-first-vite-project
type Scaffolding int

const (
	// React indicates a Vite template for a React project. This constant can be
	// used to identify if a React-specific configuration is needed.
	React Scaffolding = 1 + iota

	// ReactTs indicates a Vite template for a TypeScript React project. This
	// constant can be used to identify if a React-specific configuration is
	// needed.
	ReactTs

	// ReactSWC indicates a Vite template for a React project using SWC as the
	// compiler. This constant can be used to identify SWC-specific
	// configurations.
	ReactSwc

	// ReactSWCTs indicates a Vite template for a TypeScript React project using
	// SWC as the compiler. This constant can be used to identify SWC-specific
	// configurations for TypeScript.
	ReactSwcTs

	// Vanilla indicates a Vite template for a Vanilla JavaScript project.
	// This constant can be used to identify configurations for a basic setup
	// without frameworks.
	Vanilla

	// VanillaTs indicates a Vite template for a Vanilla TypeScript project
	// without frameworks. This constant can be used to identify
	// TypeScript-specific configurations for a basic setup.
	VanillaTs

	// Vue indicates a Vite template for a Vue.js project. This constant can be
	// used to identify if a Vue-specific configuration is needed.
	Vue

	// VueTs indicates a Vite template for a TypeScript Vue.js project. This
	// constant can be used to identify if a TypeScript Vue-specific
	// configuration is needed.
	VueTs

	// Preact indicates a Vite template for a Preact project. This constant can
	// be used to identify if a Preact-specific configuration is needed.
	Preact

	// PreactTs indicates a Vite template for a TypeScript Preact project. This
	// constant can be used to identify if a TypeScript Preact-specific
	// configuration is needed.
	PreactTs

	// Lit indicates a Vite template for a Lit project. This constant can be used
	// to identify if a Lit-specific configuration is needed.
	Lit

	// LitTs indicates a Vite template for a TypeScript Lit project. This
	// constant can be used to identify if a TypeScript Lit-specific
	// configuration is needed.
	LitTs

	// Svelte indicates a Vite template for a Svelte project. This constant can
	// be used to identify if a Svelte-specific configuration is needed.
	Svelte

	// SvelteTs indicates a Vite template for a TypeScript Svelte project. This
	// constant can be used to identify if a TypeScript Svelte-specific
	// configuration is needed.
	SvelteTs

	// Solid indicates a Vite template for a Solid project. This constant can be
	// used to identify if a Solid-specific configuration is needed.
	Solid

	// SolidTs indicates a Vite template for a TypeScript Solid project. This
	// constant can be used to identify if a TypeScript Solid-specific
	// configuration is needed.
	SolidTs

	// Qwik indicates a Vite template for a Qwik project. This constant can be
	// used to identify if a Qwik-specific configuration is needed.
	Qwik

	// QwikTs indicates a Vite template for a TypeScript Qwik project. This
	// constant can be used to identify if a TypeScript Qwik-specific
	// configuration is needed.
	QwikTs

	// None indicates that the user has opted out of using a specific
	// scaffolding. This constant can be used to specify that no template
	// configuration is desired.
	None
)

// RequiresPreamble determines if the specific scaffolding requires a
// preamble configuration.
func (s Scaffolding) RequiresPreamble() bool {
	switch s {
	case React:
		return true
	case ReactTs:
		return true
	case ReactSwc:
		return true
	case ReactSwcTs:
		return true
	default:
		return false
	}
}

// Preamble returns the preamble string associated with the Scaffolding. It
// takes a viteURL string as a parameter and returns the appropriate preamble.
func (s Scaffolding) Preamble(viteURL string) string {
	switch s {
	case React:
		return PluginReactPreamble(viteURL)
	case ReactTs:
		return PluginReactPreamble(viteURL)
	case ReactSwc:
		return PluginReactPreamble(viteURL)
	case ReactSwcTs:
		return PluginReactPreamble(viteURL)
	default:
		return ""
	}
}
