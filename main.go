package main

import "mseScraping/pdfDownloader"

func main() {
	downloader := pdfDownloader.CreateDownloader(
		"some-file-url",
		"some-file-name",
		"some-file-path",
	)

	downloader.GetPdfs()
}
