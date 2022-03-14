package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"mseScraping/pdfDownloader"
)

func main() {
	envs, err := godotenv.Read(".env")

	if err != nil {
		panic(err)
	}
	var mode string
	flag.StringVar(&mode, "mode", "", "The mode to run in, the options are download or clean.")
	flag.Parse()

	switch mode {
	case "download":
		downloader := pdfDownloader.CreateDownloader(
			envs["MSE_URL"],
			envs["RAW_PDF_PATH"],
			envs["RAW_CSV_PATH"],
			envs["ERROR_FILE_PATH"],
		)

		downloader.GetPdfs()
	default:
		log.Fatal("Please Enter the Necessary Flag (e.g. -download)", mode)
	}
}
