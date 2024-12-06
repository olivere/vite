package main

import (
	"io/fs"
	"log"
	"net/http"
)

var app AppConfig

func main() {
	app.LoadConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("/admin", adminHandler)
	mux.HandleFunc("/", appHandler)

	assetsApp, _ := fs.Sub(app.DistFS, "app")
	assetsAppFS := http.FileServerFS(assetsApp)
	mux.Handle("/assets/", assetsAppFS)

	assetsAdminApp, _ := fs.Sub(app.DistFS, "admin-app")
	assetsAdminAppFS := http.FileServerFS(assetsAdminApp)
	mux.Handle("/admin-app/assets/", http.StripPrefix("/admin-app", assetsAdminAppFS))

	log.Print("Listening on :8080...")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Data: map[string]interface{}{
			"Title": "Welcome, Admin!",
		},
	}

	renderAdminTemplate(w, "admin.gohtml", data)
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	data := &TemplateData{
		Data: map[string]interface{}{
			"Title": "Welcome, User!",
		},
	}
	renderAppTemplate(w, "app.gohtml", data)
}
