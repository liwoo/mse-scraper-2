package pdfDownloader

import (
	"fmt"
	"mseScraping/utils"
)

func CreateDownloader(
	downloadUrl string,
	pdfPath string,
	csvPath string,
	errorPath string) *PdfDownloader {

	utils.EnsureDirsExist([]string{pdfPath, csvPath, errorPath})

	return &PdfDownloader{
		DownloadUrlTemplate: downloadUrl,
		PdfPath:             pdfPath,
		CsvPath:             csvPath,
		ErrorPath:           errorPath,
	}
}

type PdfDownloader struct {
	DownloadUrlTemplate string
	PdfPath             string
	CsvPath             string
	ErrorPath           string
}

func (d *PdfDownloader) GetPdfs() {
	fmt.Println("Downloading pdfs from ", d.DownloadUrlTemplate)
}
