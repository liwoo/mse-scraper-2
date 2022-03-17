package scraper

import (
	"fmt"
	"mseScraping/utils"
)

func CreateScraper(
	downloadUrl string,
	pdfPath string,
	csvPath string,
	errorPath string,
	apiKey string,
	cleanedCSVPath string,
	cleanedJsonPath string,
	queueSize int,
	workerNum int,
	pdfStartNum string,
	pdfEndNum string) *Scraper {

	utils.EnsureDirsExist([]string{pdfPath, csvPath, errorPath, cleanedCSVPath, cleanedJsonPath})

	return &Scraper{
		DownloadUrlTemplate: downloadUrl,
		PdfPath:             pdfPath,
		CsvPath:             csvPath,
		ErrorPath:           errorPath,
		ApiKey:              apiKey,
		CleanedCSVPath:      cleanedCSVPath,
		CleanedJsonPath:     cleanedJsonPath,
		QueueSize:           queueSize,
		WorkerNum:           workerNum,
		PdfStartNum:         pdfStartNum,
		PdfEndNum:           pdfEndNum,
	}
}

type Scraper struct {
	DownloadUrlTemplate string
	PdfPath             string
	CsvPath             string
	ErrorPath           string
	ApiKey              string
	CleanedCSVPath      string
	CleanedJsonPath     string
	QueueSize           int
	WorkerNum           int
	PdfStartNum         string
	PdfEndNum           string
}

func (s *Scraper) Download() {
	fmt.Println("Downloading pdfs from ", s.DownloadUrlTemplate)
}

func (s *Scraper) Clean()  {
	fmt.Println("Saving pdfs to ", s.CleanedCSVPath)
}

func (s *Scraper) Save() {
	fmt.Println("Saving to Database from... ", s.CleanedJsonPath)
}