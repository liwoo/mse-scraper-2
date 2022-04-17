package scraper

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"mseScraping/cleaner"
	"mseScraping/downloader"
	"mseScraping/saver"
	"mseScraping/utils"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pdftables/go-pdftables-api/pkg/client"
	"github.com/sirsean/go-pool"
)

var clientCSV client.Client

func CreateScraper(
	downloadUrl string,
	pdfPath string,
	csvPath string,
	errorPath string,
	apiKey string,
	cleanedCSVPath string,
	cleanedJsonPath string,
	dbConnectionString string,
	queueSize int,
	workerNum int,
	pdfStartNum int,
	pdfEndNum int) *Scraper {

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
		DBConnectionString:  dbConnectionString,
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
	PdfStartNum         int
	PdfEndNum           int
	DBConnectionString  string
}

func (s *Scraper) Download() {
	fmt.Println("Downloading pdfs from ", s.DownloadUrlTemplate)
	var Client = http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}
	clientCSV = client.Client{
		APIKey:     s.ApiKey,
		HTTPClient: http.DefaultClient,
	}
	p := pool.NewPool(s.QueueSize, s.WorkerNum)
	p.Start()

	for i := s.PdfStartNum; i <= s.PdfEndNum; i++ {
		p.Add(downloader.MSEPdfDownloader{
			FileUrl:     fmt.Sprint(s.DownloadUrlTemplate, i),
			FileName:    fmt.Sprint(s.PdfPath, i, ".pdf"),
			FileNameCSV: fmt.Sprint(s.CsvPath, i, ".csv"),
			Client:      Client,
			CsvClient:   &clientCSV,
		})
	}

	p.Close()
}

func (s *Scraper) Clean() {
	fmt.Println("Saving pdfs to ", s.CleanedCSVPath)
	p := pool.NewPool(s.QueueSize, s.WorkerNum)
	p.Start()

	files, err := cleaner.GetAllFilesToClean(s.CsvPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		p.Add(cleaner.MSECsvCleaner{
			FileUrl:      file,
			ErrorPath:    s.ErrorPath,
			CleanCsvPath: s.CleanedCSVPath,
		})
	}

	p.Close()
}

func (s *Scraper) Save() {
	fmt.Println("Saving to Database from... ", s.CleanedJsonPath)
	p := pool.NewPool(s.QueueSize, s.WorkerNum)

	poolConfig, err := pgxpool.ParseConfig(s.DBConnectionString)
	if err != nil {
		log.Fatalln(err)
	}

	db, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalln(err)
	}

	p.Start()

	files, err := utils.GetAllCsvFiles(s.CleanedCSVPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		p.Add(saver.MSESaver{
			FileUrl:   file,
			ErrorPath: s.ErrorPath,
			Db:        db,
		})
	}

	p.Close()
}
