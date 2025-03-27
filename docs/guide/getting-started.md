---
next:
  text: 'Usage'
  link: '/guide/usage'
---

# Go-Vite Integration

This library helps you integrate a Vite frontend with a Go backend by following the [official Vite backend integration guide](https://vitejs.dev/guide/backend-integration.html).


## Overview

This library offers two approaches:

1. **Helper Function** - Generate HTML tags to link to your Vite assets
2. **HTTP Handler** - A complete `http.Handler` to serve your Vite application

## Prerequisites

Before using this library, ensure your Vite frontend is configured correctly:
::: warning IMPORTANT
- Follow the [Vite backend integration guidelines](https://vitejs.dev/guide/backend-integration.html) to configure your `vite.config.(js|ts)` file to generate the `manifest.json` for production.
:::


## Quick Start

```bash
$ go get github.com/olivere/vite
```


## Supported frameworks
This library works with all major Vite-supported frontend frameworks:
- React (with TypeScript and SWC options)
- Vue
- Vanilla JS/TS
- Preact
- Lit
- Svelte
- Solid
- Qwik


## Usage

- [Helper Function](/guide/usage#option-1-helper-function)
- [HTTP Handler](/guide/usage#option-2-http-handler)
- [Configuration](/guide/usage#configuration-options)
        


## Examples

See the [examples](/guide/examples) page.

## License

This project is licensed under the terms specified in the [LICENSE](https://raw.githubusercontent.com/olivere/vite/refs/heads/main/LICENSE) file.

