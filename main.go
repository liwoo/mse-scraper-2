package main

import (
	"mseScraping/pkg/server"
	"net/http"
)

func main() {
	app := server.Build()

	app.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to MSE Scrapper"))
	})

	app.Run()
}
