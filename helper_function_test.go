package vite_test

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/olivere/vite"
)

// from https://github.com/vitejs/vite/blob/242f550eb46c93896fca6b55495578921e29a8af/docs/guide/backend-integration.md
const exampleManifest string = `
{
  "_shared-CPdiUi_T.js": {
    "file": "assets/shared-ChJ_j-JJ.css",
    "src": "_shared-CPdiUi_T.js"
  },
  "_shared-B7PI925R.js": {
    "file": "assets/shared-B7PI925R.js",
    "name": "shared",
    "css": ["assets/shared-ChJ_j-JJ.css"]
  },
  "baz.js": {
    "file": "assets/baz-B2H3sXNv.js",
    "name": "baz",
    "src": "baz.js",
    "isDynamicEntry": true
  },
  "views/bar.js": {
    "file": "assets/bar-gkvgaI9m.js",
    "name": "bar",
    "src": "views/bar.js",
    "isEntry": true,
    "imports": ["_shared-B7PI925R.js"],
    "dynamicImports": ["baz.js"]
  },
  "views/foo.js": {
    "file": "assets/foo-BRBmoGS9.js",
    "name": "foo",
    "src": "views/foo.js",
    "isEntry": true,
    "imports": ["_shared-B7PI925R.js"],
    "css": ["assets/foo-5UjPuW-k.css"]
  }
}
`

// these are the tags we should be generating based on the manifest
const fooEntrpointTagsBlock string = `
<link rel="stylesheet" href="/assets/foo-5UjPuW-k.css">
<link rel="stylesheet" href="/assets/shared-ChJ_j-JJ.css">
<script type="module" src="/assets/foo-BRBmoGS9.js"></script>
<link rel="modulepreload" href="/assets/shared-B7PI925R.js">
`

const barEntrypointTagsBlock string = `
<link rel="stylesheet" href="/assets/shared-ChJ_j-JJ.css">
<script type="module" src="/assets/bar-gkvgaI9m.js"></script>
<link rel="modulepreload" href="/assets/shared-B7PI925R.js">
`

func getTestFS() fs.FS {
	manifestFile := fstest.MapFile{
		Data: []byte(exampleManifest),
	}
	return fstest.MapFS{
		".vite/manifest.json": &manifestFile,
	}
}

func TestFragmentContainsTagsForFooEntrpointFromManifest(t *testing.T) {
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:        getTestFS(),
		IsDev:     false,
		ViteEntry: "views/foo.js",
	})

	if err != nil {
		t.Fatal("Unable to produce Vite HTML Fragment", err)
	}

	generatedHTML := string(viteFragment.Tags)

	fooEntrypointTags := strings.Split(fooEntrpointTagsBlock, "\n")

	for _, tag := range fooEntrypointTags {
		if tag == "" {
			continue
		}

		HTMLContainsTag := strings.Contains(generatedHTML, strings.TrimSpace(tag))
		if !HTMLContainsTag {
			t.Logf(`
	------------	Generated HTML:  --- %s
			`, generatedHTML)
			t.Fatalf("Generated HTML block does not contain needed tag: %s", tag)
		}
	}
}

func TestFragmentContainsTagsForBarEntrpointFromManifest(t *testing.T) {
	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:        getTestFS(),
		IsDev:     false,
		ViteEntry: "views/bar.js",
	})

	if err != nil {
		t.Fatal("Unable to produce Vite HTML Fragment", err)
	}

	generatedHTML := string(viteFragment.Tags)

	barEntrypointTags := strings.Split(barEntrypointTagsBlock, "\n")

	for _, tag := range barEntrypointTags {
		if tag == "" {
			continue
		}

		HTMLContainsTag := strings.Contains(generatedHTML, strings.TrimSpace(tag))
		if !HTMLContainsTag {
			t.Logf(`
	------------	Generated HTML:  --- %s
			`, generatedHTML)
			t.Fatalf("Generated HTML block does not contain needed tag: %s", tag)

		}
	}
}

func TestDevModeFragmentContainsModuleTags(t *testing.T) {
	const entrypoint string = "src/main.tsx"

	viteFragment, err := vite.HTMLFragment(vite.Config{
		FS:        getTestFS(),
		IsDev:     true,
		ViteURL:   "http://localhost:5173",
		ViteEntry: entrypoint,
	})

	if err != nil {
		t.Fatal("Unable to produce Vite HTML Fragment", err)
	}

	generatedHTML := string(viteFragment.Tags)

	const viteClientTag string = `<script type="module" src="http://localhost:5173/@vite/client"></script>`
	var entrypointTag string = fmt.Sprintf(`<script type="module" src="http://localhost:5173/%s"></script>`, entrypoint)

	if !strings.Contains(generatedHTML, viteClientTag) {
		t.Fatalf("Generated HTML block does not contain: %s", viteClientTag)
	}

	if !strings.Contains(generatedHTML, entrypointTag) {
		t.Fatalf("Generated HTML block does not contain: %s", entrypointTag)
	}
}
