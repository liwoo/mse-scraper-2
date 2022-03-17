
# Malawi 🇲🇼 Stock Exchange 📈 Scraper

A scraping tool for MSE Daily Reports, which are uploaded on [their site](https://mse.co.mw/index.php?route=market/market/report).  The tool basically goes through different daily reports (according tothe configuration) and downloads the PDFs before converting them to CSV and then eventually uploading to a Postgres Database.  What you do with the data is up to you! Just make sure you have a VPN when scraping.

We use [PDF Tables](https://pdftables.com/) to convert the PDFs to CSV, but you may use your own converter, at which point some of the logic in the `downloader` will need to be tweaked.
## Authors

- [@liwoo](https://www.github.com/liwoo)
- [@sevenreup](https://www.github.com/sevenreup)

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

- `MSE_URL` the absolute URL where a PDF is to be found, without the PDF Number
- `PDF_START_NO` PDF Start Number
- `PDF_END_NO` PDF End Number
- `RAW_PDF_PATH` Relative Project Path where you want to save PDFs
- `RAW_CSV_PATH` Relative Project Path where you want to save uncleaded CSVs
- `ERROR_FILE_PATH` Relative Project Path where you want to save Errors
- `CLEANED_CSV_PATH` Relative Project Path where you want to save cleaned CSVs
- `QUEUE_SIZE` Pool Maximum Queue Size
- `WORKER_NUM` Pool Number of workers
- `PDFTABLES_API_KEY` PDF Tables API Key


## Installation

After cloning this repo, make sure you change copy the `example.env` into an `.env` and replace all the values in there with sensible configurations.

You may then build and run the program with the following flags

```bash
  go build -o scraper
  
  ./scraper -mode download
  
  #wait for completion

  ./scraper -mode clean

  #wait for completion

  ./scraper -mode save
  
  #wait for completion
```

Any errors incurred will be both logged in the terminal as well as recorded in the error path you provide.  You may handle the errors however you see fit - including manually converting and saving them.

## Acknowledgements

- [PDF Tables](https://pdftables.com/)