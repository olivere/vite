# Vite Multi-App Example

This repository demonstrates how to set up and serve two separate Vite applications integrated into a Go backend.

## Overview

We have two Vite-based frontend applications:

1. **Admin App**: Located in `frontend/admin-app` (Vite + ReactTS).
2. **Main App**: Located in `frontend/app` (Vite + Vue.js).

Each application is served with its own Vite development server during development and shares a common Go backend in production.

## Setting Up Vite Tags in Go

To generate Vite tags for embedding the frontend assets, we use two separate functions: `ViteAppTags` and `ViteAdminTags`.

### Example Configuration for `ViteAppTags`

This function generates tags for the **Main App**:

```go
func ViteAppTags() template.HTML {
	var err error
	assetsApp, _ := fs.Sub(app.DistFS, "app")
	viteConfig := vite.Config{
		IsDev:        app.InDevelopment,
		ViteURL:      "http://localhost:5173", // Development URL for 'app'
		ViteEntry:    "src/main.tsx",
		ViteTemplate: vite.ReactTs,
		FS:           assetsApp,
	}

	viteFragment, err := vite.HTMLFragment(viteConfig)
	if err != nil {
		log.Printf("Vite fragment error: %v", err)
		return "<!-- Vite fragment error -->"
	}

	return viteFragment.Tags
}
```

### Example Configuration for `ViteAdminTags`

This function generates tags for the **Admin App**:

```go
func ViteAdminTags() template.HTML {
	var err error
	assetsApp, _ := fs.Sub(app.DistFS, "admin-app")
	viteConfig := vite.Config{
		IsDev:           app.InDevelopment,
		ViteURL:         "http://localhost:5174/admin", // Development URL for 'admin'
		ViteEntry:       "src/main.js",
		ViteTemplate:    vite.Vue,
		FS:              assetsApp,
		AssetsURLPrefix: "/admin-app", // Custom prefix
	}

	viteFragment, err := vite.HTMLFragment(viteConfig)
	if err != nil {
		log.Printf("Vite fragment error: %v", err)
		return "<!-- Vite fragment error -->"
	}

	return viteFragment.Tags
}
```

### Key Notes

1. **Different Vite URLs**: Each application has a unique `ViteURL` during development:
   - Main App: `http://localhost:5173`
   - Admin App: `http://localhost:5174/admin`

2. **Custom `AssetsURLPrefix`**: The `ViteAdminTags` function includes an `AssetsURLPrefix` for the **Admin App**, which modifies asset paths in production. For example:
   ```html
   <script type="module" src="/admin-app/assets/main-DVg2CzpX.js"></script>
   <link rel="modulepreload" href="/admin-app/assets/main-DVg2CzpX.js">
   ```

3. **Serving Assets in Production**: Custom prefixes require a handler to serve the assets. In `main.go`, this is achieved as follows:

   ```go
   assetsAdminApp, _ := fs.Sub(app.DistFS, "admin-app")
   assetsAdminAppFS := http.FileServer(http.FS(assetsAdminApp))
   mux.Handle("/admin-app/assets/", http.StripPrefix("/admin-app", assetsAdminAppFS))
   ```

   Here, `app.DistFS` points to the directory containing the `dist` files. You can use any file system, such as `os.DirFS` or `embed`.

## Running the Application

### Development Mode

To run the application in development:

1. Start all services:
   ```sh
   make start-all
   ```

2. Open the browser:
   - Main App: [http://localhost:8080/](http://localhost:8080/) → `Welcome, User!`
   - Admin App: [http://localhost:8080/admin](http://localhost:8080/admin) → `Welcome, Admin!`

### Production Mode

To build and run the application in production:

1. Build the assets:
   ```sh
   make build-all
   ```

2. Run the backend:
   ```sh
   APP_ENV=production ./backend
   ```

3. Open the browser:
   - Main App: [http://localhost:8080/](http://localhost:8080/)
   - Admin App: [http://localhost:8080/admin](http://localhost:8080/admin)

## Additional Notes

- The `make start-all` and `make build-all` commands manage both Vite apps and the backend.
- Hot Module Reloading (HMR) is supported during development for both apps.
- In production, static files are served directly from the Go binary or filesystem.

Enjoy your multi-app setup with Vite and Go!