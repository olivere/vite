package main

import (
	"io/fs"
	"os"
)

type AppConfig struct {
	DistFS        fs.FS
	Environment   string
	InProduction  bool
	InStaging     bool
	InDevelopment bool
}

func (app *AppConfig) LoadConfig() error {
	app.Environment = os.Getenv("APP_ENV")
	app.InProduction = app.Environment == "production"
	app.InStaging = app.Environment == "staging"
	app.InDevelopment = app.Environment == ""
	app.DistFS = os.DirFS("dist")
	return nil
}
