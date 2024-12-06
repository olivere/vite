package main

import (
	"html/template"
	"io/fs"
	"log"
	"net/http"

	"github.com/olivere/vite"
)

type TemplateData struct {
	Data map[string]interface{}
}

func renderTemplate(w http.ResponseWriter, tmpl string, funcs template.FuncMap, data *TemplateData) {
	t, err := template.New(tmpl).Funcs(funcs).ParseFiles(tmpl)
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		log.Println("Template parsing error:", err)
		return
	}

	err = t.Execute(w, data.Data)
	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		log.Println("Template execution error:", err)
	}
}

func renderAppTemplate(w http.ResponseWriter, tmpl string, data *TemplateData) {
	tmplFuncs := template.FuncMap{
		"viteTags": ViteAppTags,
	}

	renderTemplate(w, tmpl, tmplFuncs, data)
}

func renderAdminTemplate(w http.ResponseWriter, tmpl string, data *TemplateData) {
	tmplFuncs := template.FuncMap{
		"viteTags": ViteAdminTags,
	}

	renderTemplate(w, tmpl, tmplFuncs, data)
}

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
