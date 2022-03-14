package pdfDownloader

import "net/http"

type MseFileDownload struct {
	downloadUrl string
	pdfFileName string
	csvFileName string
	client      http.Client
}
