package downloader

import (
	"fmt"
	"net/http"
)

type MSEPdfDownloader struct {
	FileUrl     string
	FileName    string
	FileNameCSV string
	Client      http.Client
}

func nonStandardDailyReport(size int64) bool {
	return size < 49000 || size > 80000
}

func (m *MSEPdfDownloader) Perform() {
	fmt.Println("Downloading MSE PDF file...")
}

func (m MSEPdfDownloader) SavePDF() (int64, error) {
	//TODO: Implement
	return 1, nil
}

func (m MSEPdfDownloader) ConvertToCSV() error  {
	//TODO: Implement
	return nil
}
