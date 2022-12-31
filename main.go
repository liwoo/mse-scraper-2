package main

import (
	"mseScraping/api/clean"
	"mseScraping/api/scrapper"
	"mseScraping/pkg/server"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	app := server.Build()

	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to MSE Scrapper"))
	})

	app.Router.Route("/scrapper", func(r chi.Router) {
		r.Get("/download", scrapper.DownloadRange)
		r.Get("/daily", scrapper.DownloadDaily)
	})

	// Clean Date
	app.Router.Route("/clean", func(r chi.Router) {
		r.Get("/", clean.CleanDatabase)
	})

	app.Run()
}
