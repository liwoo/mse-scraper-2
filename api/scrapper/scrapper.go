package scrapper

import "net/http"

func DownloadRange(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Download Range"))
}

func DownloadDaily(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Download Daily"))
}