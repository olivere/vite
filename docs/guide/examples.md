---
prev:
  text: 'Usage'
  link: '/guide/usage'

next: false  
---

# Examples

This section provides various examples demonstrating how to use this library with different setups. Each example can be quickly tried using the `npx degit` command provided.

## Simple Helper Function

::: tip Minimal Setup
A minimal React app with Go backend integration using basic helper functions.
:::

```bash
$ npx degit olivere/vite/examples/helper-function-basic my-helper-app
$ cd my-helper-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/olivere/vite/tree/main/examples/helper-function-basic)

## Basic with Handler

::: tip Standard Approach
Demonstrates a basic React app with Go backend using the standard handler approach.
:::

```bash
$ npx degit olivere/vite/examples/basic my-basic-app
$ cd my-basic-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/olivere/vite/tree/main/examples/basic)


## Multi-Page Application

::: tip Multiple Entry Points
For Vite apps with multiple entry points, demonstrating how to create separate handlers with the `ViteEntry` field.
:::

```bash
$ npx degit olivere/vite/examples/multi-page-app my-multi-page-app
$ cd my-multi-page-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/olivere/vite/tree/main/examples/multi-page-app)



## Multiple Vite Instances Application

::: tip Multiple Vite instances
For Apps where you need to manage different vite instances, like admin panel(Vite + ReactTS) and main app (Vite + Vue)
:::

```bash
$ npx degit olivere/vite/examples/multiple-vite-apps multiple-vite-apps
$ cd my-multi-page-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/olivere/vite/tree/main/examples/multiple-vite-apps)

## Template Registration

::: tip Custom Templates
Shows how to use custom HTML templates in your Go backend for serving different React pages.
:::

```bash
$ npx degit olivere/vite/examples/template-registry my-template-app
$ cd my-template-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/olivere/vite/tree/main/examples/template-registry)

## Inertia.js Integration

::: warning Third-Party Integration
Example of using Golang with `net/http`, Inertia.js and this library for managing Vite assets.
:::

```bash
$ npx degit danclaytondev/go-inertia-vite my-inertia-app
$ cd my-inertia-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/danclaytondev/go-inertia-vite)


## Router Application

::: warning Advanced Implementation
A more complex application with a Go backend serving a Vite-based app using TanStack Router and TanStack Query libraries.
:::

```bash
$ npx degit olivere/vite/examples/router my-router-app
$ cd my-router-app
# Follow setup instructions in the example README
```

[View code →](https://github.com/olivere/vite/tree/main/examples/router)