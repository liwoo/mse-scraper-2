package main

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	scraper2 "mseScraping/scraper"
	"strconv"
)

func main() {
	envs, err := godotenv.Read(".env")

	if err != nil {
		panic(err)
	}
	var mode string
	flag.StringVar(&mode, "mode", "", "The mode to run in, the options are download or clean.")
	flag.Parse()

	queueSize, queueSizeError := strconv.Atoi(envs["QUEUE_SIZE"])
	workerNum, workerNumError := strconv.Atoi(envs["WORKER_NUM"])

	if queueSizeError != nil || workerNumError != nil {
		log.Panicf("Cannot convert Queue Size or Worker Nums to Int")
	}

	scraper := scraper2.CreateScraper(
		envs["MSE_URL"],
		envs["RAW_PDF_PATH"],
		envs["RAW_CSV_PATH"],
		envs["ERROR_FILE_PATH"],
		envs["PDFTABLES_API_KEY"],
		envs["CLEANED_CSV_PATH"],
		envs["CLEANED_JSON_PATH"],
		queueSize,
		workerNum,
		envs["PDF_START_NO"],
		envs["PDF_END_NO"])

	switch mode {
	case "download":
		scraper.Download()
	case "clean":
		scraper.Clean()
	case "save":
		scraper.Save()
	default:
		log.Fatal("Please Enter the Necessary Flag (e.g. -download)", mode)
	}
}
