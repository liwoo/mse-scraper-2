package server

import (
	"fmt"
	"mseScraping/pkg/conf"
	"mseScraping/utils"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type HandleWithBody func(body any)

type Server struct {
	Conf   conf.Conf
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

	initPaths(&config)

	r := chi.NewRouter()
	app := Server{
		Router: r,
		Conf:   config,
	}
	return app
}

func initPaths(conf *conf.Conf) {
	utils.EnsureDirsExist([]string{conf.PdfPath, conf.CleanedCSVPath, conf.ErrorPath, conf.CsvPath, conf.CleanedJsonPath})
}

func (app Server) Run() {
	http.ListenAndServe(":"+app.Conf.Port, app.Router)
}
