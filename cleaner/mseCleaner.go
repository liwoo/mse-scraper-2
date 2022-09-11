package cleaner

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type DailyCompanyRate struct {
	NO         string  `json:"no" validate:"string,required=true"`
	HIGH       float64 `json:"high" validate:"number,required=false"`
	LOW        float64 `json:"low" validate:"number,required=false"` // required true for everything else
	CODE       string  `json:"code" validate:"string,required=true"`
	BUY        float64 `json:"buy" validate:"number,required=false"`
	SELL       float64 `json:"sell" validate:"number,required=false"`
	PCP        float64 `json:"pcp" validate:"number,required=true"`
	TCP        float64 `json:"tcp" validate:"number,required=true"`
	VOL        float64 `json:"vol" validate:"number,required=true"`
	DIVNET     float64 `json:"div_net" validate:"number,required=false"`
	DIVYIELD   float64 `json:"div_yield" validate:"number,required=false"`
	EEARNYIELD float64 `json:"earn_yield" validate:"number,required=true"`
	PERATIO    float64 `json:"pe_ratio" validate:"number,required=true"`
	PBVRATION  float64 `json:"pbv_ratio" validate:"number,required=true"`
	CAP        float64 `json:"cap" validate:"number,required=true"`
	PROFIT     float64 `json:"profit" validate:"number,required=true"`
	SHARES     float64 `json:"shares" validate:"number,required=true"`
}

type CleanedData struct {
	dailyRates []DailyCompanyRate
	date       string
	errors     []string
}

type MSECsvCleaner struct {
	FileUrl      string
	ErrorPath    string
	CleanCsvPath string
}

func GetAllFilesToClean(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) != ".csv" {
			return nil
		}
		files = append(files, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func (u MSECsvCleaner) Perform() {
	data, err := Clean(u.FileUrl, u.ErrorPath, u.CleanCsvPath)
	if err != nil {
		log.Println(err)
	} else {
		if len(data.errors) > 0 {
			log.Println("Cleaner has errors, ", data.errors)
		}
	}
}

func Clean(csvFile string, errorPath string, cleanCSVPath string) (*CleanedData, error) {
	fmt.Println("Cleaning file: ", csvFile)
	var rates []DailyCompanyRate
	var cleaningErrors []string
	var date string
	docNum := getDocName(csvFile)

	fileBytes, err := os.ReadFile(csvFile)

	if err != nil {
		return nil, err
	}

	dateRegex, err := regexp.Compile("As On Date:.*")
	if err != nil {
		return nil, err
	}

	for _, match := range dateRegex.FindAllString(string(fileBytes), -1) {
		d, err := GetDate(match, docNum)
		if err != nil {
			cleaningErrors = append(cleaningErrors, err.Error())
		}
		date = d
	}

	dataRegex, err := regexp.Compile("Daily(?s)(.*)(?:Indices)")
	if err != nil {
		return nil, errors.New(fmt.Sprint(csvFile, err))
	}

	r := csv.NewReader(bytes.NewBuffer(dataRegex.Find(fileBytes)))
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	records, err := r.ReadAll()

	if err != nil {
		return nil, errors.New(fmt.Sprint(csvFile, err))
	}
	var rate DailyCompanyRate
	for i, word := range records {
		if isInt(word[0]) {
			if len(word) == 17 {
				rate.NO = strings.TrimSpace(word[0])
				rate.HIGH = parseFloat(strings.TrimSpace(word[1]))
				rate.LOW = parseFloat(strings.TrimSpace(word[2]))
				rate.CODE = strings.TrimSpace(word[3])
				rate.BUY = parseFloat(strings.TrimSpace(word[4]))
				rate.SELL = parseFloat(strings.TrimSpace(word[5]))
				rate.PCP = parseFloat(strings.TrimSpace(word[6]))
				rate.TCP = parseFloat(strings.TrimSpace(word[7]))
				rate.VOL = parseFloat(strings.TrimSpace(word[8]))
				rate.DIVNET = parseFloat(strings.TrimSpace(word[9]))
				rate.DIVYIELD = parseFloat(strings.TrimSpace(word[10]))
				rate.EEARNYIELD = parseFloat(strings.TrimSpace(word[11]))
				rate.PERATIO = parseFloat(strings.TrimSpace(word[12]))
				rate.PBVRATION = parseFloat(strings.TrimSpace(word[13]))
				rate.CAP = parseFloat(strings.TrimSpace(word[14]))
				rate.PROFIT = parseFloat(strings.TrimSpace(word[15]))
				rate.SHARES = parseFloat(strings.TrimSpace(word[16]))

				_, err := Verify(rate)
				if err != nil {
					cleaningErrors = append(cleaningErrors, fmt.Sprintf("line: %d, number: %q {\n %s \n}", i, word[0], err.Error()))
				}

				rates = append(rates, rate)
			}
		}
	}

	if len(cleaningErrors) > 0 {
		affected := fmt.Sprintf("File: %q", csvFile)
		cleaningErrors = append([]string{affected}, cleaningErrors...)
		err := logErrors(cleaningErrors, errorPath, date)
		if err != nil {
			fmt.Println(csvFile, err)
		}
	}

	err = func() error {
		file, err := os.Create(fmt.Sprintf("%s%s.csv", cleanCSVPath, strings.ReplaceAll(date, "/", "-")))
		if err != nil {
			return err
		}
		defer file.Close()
		w := csv.NewWriter(file)
		for _, rate := range rates {
			if err := w.Write([]string{
				rate.NO, fmt.Sprint(rate.HIGH), fmt.Sprint(rate.LOW), rate.CODE,
				fmt.Sprintf("%.2f", rate.BUY), fmt.Sprintf("%.2f", rate.SELL), fmt.Sprintf("%.2f", rate.PCP), fmt.Sprintf("%.2f", rate.TCP), fmt.Sprintf("%.2f", rate.VOL),
				fmt.Sprint(rate.DIVNET), fmt.Sprint(rate.DIVYIELD), fmt.Sprintf("%.2f", rate.EEARNYIELD),
				fmt.Sprintf("%.2f", rate.PERATIO), fmt.Sprintf("%.2f", rate.PBVRATION), fmt.Sprintf("%.2f", rate.CAP), fmt.Sprintf("%.2f", rate.PROFIT), fmt.Sprintf("%.2f", rate.SHARES)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		w.Flush()

		if err := w.Error(); err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		return nil, errors.New(fmt.Sprint(csvFile, err))
	}

	fmt.Println("Finished file: ", csvFile)

	return &CleanedData{
		dailyRates: rates,
		date:       date,
		errors:     cleaningErrors,
	}, nil
}

func logErrors(errors []string, path string, date string) error {
	filePath := fmt.Sprintf("%s%s-error.txt", path, strings.ReplaceAll(date, "/", "-"))
	fmt.Println(filePath, path, date)
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err2 := f.WriteString(strings.Join(errors, "\n"))
	if err2 != nil {
		return err
	}

	return nil
}

func parseFloat(value string) float64 {
	if strings.Contains(value, ",") {
		value = strings.ReplaceAll(value, ",", "")
	}
	if strings.Contains(value, "(") {
		raw := strings.ReplaceAll(value, "(", "")
		raw = strings.ReplaceAll(raw, ")", "")

		raw = strings.TrimSpace(raw)
		value = fmt.Sprint("-", raw)
	}

	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return -1
	}
	return f
}

func isInt(value string) bool {
	_, err := strconv.ParseInt(value, 10, 8)
	return err == nil
}

func Verify(rate DailyCompanyRate) (bool, error) {
	err := validateStruct(rate)
	if len(err) > 0 {
		var listOfErrors []string
		for _, x := range err {
			listOfErrors = append(listOfErrors, x.Error())
		}
		return false, errors.New(strings.Join(listOfErrors, "\n"))
	}
	return true, nil
}

func isEmpty(value string) bool {
	return len(value) <= 0
}

func GetDate(line string, docNum string) (string, error) {
	r, _ := regexp.Compile("\\d?\\d/\\d\\d/\\d\\d\\d\\d")
	if r.Match([]byte(line)) {
		match := r.FindString(line)

		if isEmpty(match) {
			return docNum, errors.New("could not find date match in string")
		}
		t, err := time.Parse(checkDateFormat(match), match)
		if err != nil {
			fmt.Println(err)
			return docNum, errors.New("failed to parse date")
		}
		return t.Format("2006-01-02"), nil
	} else {
		return docNum, fmt.Errorf("line does not contain date, : %s", line)
	}
}

func checkDateFormat(date string) string {
	segs := strings.Split(date, "/")
	if len(segs[0]) == 1 {
		return "2/01/2006"
	} else {
		return "02/01/2006"
	}
}

func getDocName(fileName string) string {
	return filepath.Base(fileName)
}
