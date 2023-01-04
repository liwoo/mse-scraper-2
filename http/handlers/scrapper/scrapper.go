package scrapper

import (
	"fmt"
	"log"
	"mseScraping/http/handlers/scrapper/cleaner"
	"mseScraping/http/handlers/scrapper/downloader"
	"mseScraping/http/handlers/scrapper/saver"
	"mseScraping/pkg/conf"
	"mseScraping/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"database/sql"

	"github.com/gocolly/colly"
	"github.com/pdftables/go-pdftables-api/pkg/client"
	"github.com/sirsean/go-pool"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var clientCSV client.Client

func DownloadRange(w http.ResponseWriter, r *http.Request, conf conf.Conf) {
	start, err := strconv.Atoi(r.URL.Query().Get("start"))
	end, err1 := strconv.Atoi(r.URL.Query().Get("end"))

	if err != nil || err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Range failed to parse"))
		return
	}

	scrap(conf, start, end)
	w.Write([]byte("Download Range"))
}

func DownloadDaily(w http.ResponseWriter, r *http.Request, conf conf.Conf) {
	date := time.Now()
	c := colly.NewCollector()
	selector := fmt.Sprintf("td:contains(\"Daily %v %v %v\")", date.Day(), date.Month(), date.Year())

	c.OnHTML(selector, func(e *colly.HTMLElement) {
		link := e.DOM.Parent().Find("a")
		url, exists := link.Attr("href")
		if !exists {
			return
		}
		segs := strings.Split(url, "/")
		numRaw := segs[len(segs)-1]
		currentNumber, err := strconv.ParseInt(numRaw, 10, 32)

		if err != nil {
			return
		}
		fmt.Println("Daily report found: ", currentNumber)
		scrap(conf, int(currentNumber), int(currentNumber))
	})

	c.Visit(conf.DownloadUrlTemplate)
	w.Write([]byte("Download Daily"))
}

func scrap(conf conf.Conf, start int, end int) {
	currenntBatchPath := fmt.Sprintf("%v_%v_%v/", time.Now().Format("2006_01_02"), start, end)
	downloadPath := download(conf, start, end, currenntBatchPath)
	cleanedPath := clean(conf, currenntBatchPath, downloadPath)
	save(conf, currenntBatchPath, cleanedPath)
}

func download(s conf.Conf, start int, end int, currentBatchPath string) string {
	fmt.Println("Downloading pdfs from ", s.DownloadUrlTemplate)
	csvDownloadPath := fmt.Sprint(s.CsvPath, currentBatchPath)
	pdfDownloadPath := fmt.Sprint(s.PdfPath, currentBatchPath)
	utils.EnsureDirsExist([]string{csvDownloadPath, pdfDownloadPath})
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

	for i := start; i <= end; i++ {
		p.Add(downloader.MSEPdfDownloader{
			FileUrl:     fmt.Sprint(s.DownloadUrlTemplate, i),
			FileName:    fmt.Sprint(pdfDownloadPath, i, ".pdf"),
			FileNameCSV: fmt.Sprint(csvDownloadPath, i, ".csv"),
			Client:      Client,
			CsvClient:   &clientCSV,
		})
	}

	p.Close()

	return csvDownloadPath
}

func clean(s conf.Conf, currentBatchPath string, downloadPath string) string {
	cleanPath := fmt.Sprint(s.CleanedCSVPath, currentBatchPath)
	fmt.Println("Cleaning csvs to ", cleanPath)
	utils.EnsureDirsExist([]string{fmt.Sprint(s.ErrorPath, currentBatchPath), cleanPath})
	p := pool.NewPool(s.QueueSize, s.WorkerNum)
	p.Start()

	files, err := cleaner.GetAllFilesToClean(downloadPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		p.Add(cleaner.MSECsvCleaner{
			FileUrl:      file,
			ErrorPath:    fmt.Sprint(s.ErrorPath, currentBatchPath),
			CleanCsvPath: cleanPath,
		})
	}

	p.Close()

	return cleanPath
}

func save(s conf.Conf, currentBatchPath string, cleanedCsvPath string) {
	cleanedJson := fmt.Sprint(s.CleanedJsonPath, currentBatchPath)
	fmt.Println("Saving to Database from... ", cleanedCsvPath)
	utils.EnsureDirsExist([]string{cleanedJson, fmt.Sprint(s.ErrorPath, currentBatchPath)})
	p := pool.NewPool(s.QueueSize, s.WorkerNum)

	var pgconn *pgdriver.Connector = pgdriver.NewConnector(pgdriver.WithDSN(s.DBConnectionString))
	// pgconn.Config().TLSConfig = nil
	psdb := sql.OpenDB(pgconn)
	db := bun.NewDB(psdb, pgdialect.New())

	p.Start()

	files, err := utils.GetAllCsvFiles(cleanedCsvPath)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		p.Add(saver.MSESaver{
			FileUrl:   file,
			ErrorPath: fmt.Sprint(s.ErrorPath, currentBatchPath),
			Db:        db,
		})
	}

	p.Close()
}
