# `dist` Directory

This directory is used to store the built output from the Vite build process.

In development mode, this directory might be empty if the build process has not
yet been run. The Go application requires at least one file to be present in
this directory to embed its contents using `go:embed`. This README file ensures
the directory is not empty.
