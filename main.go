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

	app.Router.Route("/api/scrapper", func(r chi.Router) {
		r.Get("/download", func(w http.ResponseWriter, r *http.Request) {
			scrapper.DownloadRange(w, r, app.Conf)
		})
		r.Get("/daily", func(w http.ResponseWriter, r *http.Request) {
			scrapper.DownloadDaily(w, r, app.Conf)
		})
	})

	// Clean Date
	app.Router.Route("/api/clean", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			clean.CleanDatabase(w, r, app.Conf)
		})
	})

	app.Run()
}
