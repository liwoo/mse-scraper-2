package server

import (
	"fmt"
	"mseScraping/pkg/conf"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type HandleWithBody func(body any)

type Server struct {
	conf   conf.Conf
	Router *chi.Mux
}

func Build() Server {
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development"
	}
	godotenv.Load(".env." + environment + ".local")
	godotenv.Load(".env." + environment)
	godotenv.Load()

	config := conf.Conf{}
	err := conf.Parse(&config)

	if err != nil {
		fmt.Println(err)
	}

	r := chi.NewRouter()
	app := Server{
		Router: r,
		conf:   config,
	}
	return app
}

func (app Server) Run() {
	http.ListenAndServe(":"+app.conf.PORT, app.Router)
}
