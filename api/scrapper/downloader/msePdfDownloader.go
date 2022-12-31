package downloader

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/pdftables/go-pdftables-api/pkg/client"
)

type MSEPdfDownloader struct {
	FileUrl     string
	FileName    string
	FileNameCSV string
	Client      http.Client
	CsvClient   *client.Client
}

func nonStandardDailyReport(size int64) bool {
	//32000 - 36000 are the new format
	//14000 - 16000 are the old format
	//TODO: We might need to check if we are dealing with
	//the old or new format
	return size < 32000 || size > 200000
}

func (m MSEPdfDownloader) Perform() {

	size, err := m.SavePDF()

	if err != nil {
		log.Fatal(err)
	}

	err = m.ConvertToCSV(size)

	if err != nil {
		log.Fatal(err)
	}

}

func (m MSEPdfDownloader) SavePDF() (int64, error) {
	fmt.Println("Downloading MSE PDF file... ", m.FileUrl)
	file, err := os.Create(m.FileName)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	resp, err := m.Client.Get(m.FileUrl)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	size, err := io.Copy(file, resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Downloaded a file %s with size %d\n", m.FileName, size)
	return size, nil
}

func (m MSEPdfDownloader) ConvertToCSV(size int64) error {
	fmt.Println("Converting MSE PDF file to CSV ", m.FileName)
	file, err := os.Open(m.FileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {

		cerr := file.Close()
		if cerr == nil {
			err = cerr
		}

		if nonStandardDailyReport(size) {
			fmt.Printf("Unexcepted PDF's size. File %s,Size  %d\n", m.FileName, size)
			e := os.Remove(m.FileName)
			if e != nil {
				log.Fatal(e)
			}
		}

		if !nonStandardDailyReport(size) {
			err := attemptCsvConversion(m)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Done converting %s to CSV\n", m.FileName)
		}
	}()

	return nil
}

func attemptCsvConversion(m MSEPdfDownloader) error {
	file, err := os.Open(m.FileName)

	if err != nil {
		return err
	}

	defer file.Close()

	csvFile, err := os.Create(m.FileNameCSV)

	if err != nil {
		return err
	}

	defer csvFile.Close()

	converted, err := m.CsvClient.Do(file, client.FormatCSV)

	if err != nil {
		return err
	}

	_, err = io.Copy(csvFile, converted)

	if err != nil {
		return err
	}
	return nil
}
