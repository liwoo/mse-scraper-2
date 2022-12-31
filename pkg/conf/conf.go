package conf

type Conf struct {
	Port                string `env:"PORT"`
	DownloadUrlTemplate string `env:"MSE_URL"`
	PdfPath             string `env:"RAW_PDF_PATH"`
	CsvPath             string `env:"RAW_CSV_PATH"`
	ErrorPath           string `env:"ERROR_FILE_PATH"`
	ApiKey              string `env:"PDFTABLES_API_KEY"`
	CleanedCSVPath      string `env:"CLEANED_CSV_PATH"`
	CleanedJsonPath     string `env:"CLEANED_JSON_PATH"`
	QueueSize           int    `env:"QUEUE_SIZE"`
	WorkerNum           int    `env:"WORKER_NUM"`
	DBConnectionString  string `env:"DB_CONNECTION_STRING"`
}
