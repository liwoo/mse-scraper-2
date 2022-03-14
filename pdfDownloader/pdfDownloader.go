package pdfDownloader

import "fmt"

func CreateDownloader(fileUrl string, fileName string, fileNameCsv string) *PdfDownloader {
	return &PdfDownloader{
		FileUrl:     fileUrl,
		FileName:    fileName,
		FileNameCsv: fileNameCsv,
	}
}

type PdfDownloader struct {
	FileUrl     string
	FileName    string
	FileNameCsv string
}

func (d *PdfDownloader) GetPdfs() {
	fmt.Println("Downloading pdfs...")
}
