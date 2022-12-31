package scrapper

import (
	"fmt"
	"log"
	"mseScraping/api/scrapper/cleaner"
	"mseScraping/api/scrapper/downloader"
	"mseScraping/api/scrapper/saver"
	"mseScraping/pkg/conf"
	"mseScraping/utils"
	"net/http"

	"database/sql"

	"github.com/pdftables/go-pdftables-api/pkg/client"
	"github.com/sirsean/go-pool"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var clientCSV client.Client

func DownloadRange(w http.ResponseWriter, r *http.Request, conf conf.Conf) {
	download(conf)
	clean(conf)
	save(conf)
	w.Write([]byte("Download Range"))
}

func DownloadDaily(w http.ResponseWriter, r *http.Request, conf conf.Conf) {
	w.Write([]byte("Download Daily"))
}

func download(s conf.Conf) {
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

func clean(s conf.Conf) {
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

func save(s conf.Conf) {
	fmt.Println("Saving to Database from... ", s.CleanedJsonPath)
	p := pool.NewPool(s.QueueSize, s.WorkerNum)

	var pgconn *pgdriver.Connector = pgdriver.NewConnector(pgdriver.WithDSN(s.DBConnectionString))
	// pgconn.Config().TLSConfig = nil
	psdb := sql.OpenDB(pgconn)
	db := bun.NewDB(psdb, pgdialect.New())

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
