package cleaner

import (
	"bufio"
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
	NO         string  `json:"no"`
	HIGH       string  `json:"high"`
	LOW        string  `json:"low"`
	CODE       string  `json:"code"`
	BUY        float64 `json:"buy"`
	SELL       float64 `json:"sell"`
	PCP        float64 `json:"pcp"`
	TCP        float64 `json:"tcp"`
	VOL        float64 `json:"vol"`
	DIVNET     float64 `json:"div_net"`
	DIVYIELD   float64 `json:"div_yield"`
	EEARNYIELD float64 `json:"earn_yield"`
	PERATIO    float64 `json:"pe_ratio"`
	PBVRATION  float64 `json:"pbv_ratio"`
	CAP        float64 `json:"cap"`
	PROFIT     float64 `json:"profit"`
	SHARES     float64 `json:"shares"`
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
	data := Clean(u.FileUrl, u.ErrorPath, u.CleanCsvPath)
	if len(data.errors) > 0 {
		log.Fatalln("Cleaner has errors, ", data.errors)
	}
}

func Clean(csvFile string, errorPath string, cleanCSVPath string) CleanedData {
	fmt.Println("Cleaning file: ", csvFile)
	var csvRaw []string
	var rates []DailyCompanyRate
	var errors []string
	var date string
	docNum := getDocName(csvFile)

	file, err := os.Open(csvFile)

	if err != nil {
		log.Fatalf("failed to open file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var i int16 = 0

	for scanner.Scan() {
		//do work
		line := scanner.Text()
		if i == 8 {
			d, err := GetDate(line, docNum)
			if err != nil {
				errors = append(errors, err.Error())
			}
			date = d
		}

		if i > 26 && i < 43 {
			csvRaw = append(csvRaw, line)
		}

		i += 1
	}

	b := []byte(strings.Join(csvRaw, "\n"))

	r := csv.NewReader(bytes.NewBuffer(b))
	records, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	var rate DailyCompanyRate
	for _, word := range records {
		if len(word) == 17 {
			_, err := Verify(word)
			if err != nil {
				errors = append(errors, err.Error())
			}
			rate.NO = strings.TrimSpace(word[0])
			rate.HIGH = strings.TrimSpace(word[1])
			rate.LOW = strings.TrimSpace(word[2])
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
			rates = append(rates, rate)
		}
	}

	if len(errors) > 0 {
		affected := fmt.Sprintf("File: %q", csvFile)
		errors = append([]string{affected}, errors...)
		logErrors(errors, errorPath, date)
	}

	func() {
		file, err := os.Create(fmt.Sprintf("%s%s.csv", cleanCSVPath, strings.ReplaceAll(date, "/", "-")))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		w := csv.NewWriter(file)
		for _, rate := range rates {
			if err := w.Write([]string{
				rate.NO, rate.HIGH, rate.LOW, rate.CODE,
				fmt.Sprint(rate.BUY), fmt.Sprint(rate.SELL), fmt.Sprint(rate.PCP), fmt.Sprint(rate.TCP), fmt.Sprint(rate.VOL),
				fmt.Sprint(rate.DIVNET), fmt.Sprint(rate.DIVYIELD), fmt.Sprint(rate.EEARNYIELD),
				fmt.Sprint(rate.PERATIO), fmt.Sprint(rate.PBVRATION), fmt.Sprint(rate.CAP), fmt.Sprint(rate.PROFIT), fmt.Sprint(rate.SHARES)}); err != nil {
				log.Fatalln("error writing record to csv:", err)
			}
		}

		w.Flush()

		if err := w.Error(); err != nil {
			log.Fatal(err)
		}
	}()

	fmt.Println("Finished file: ", csvFile)

	return CleanedData{
		dailyRates: rates,
		date:       date,
		errors:     errors,
	}
}

func logErrors(errors []string, path string, date string) {
	f, err := os.Create(fmt.Sprintf("%s%s-error.txt", path, strings.ReplaceAll(date, "/", "-")))
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	_, err2 := f.WriteString(strings.Join(errors, "\n"))
	if err2 != nil {
		log.Fatal(err)
	}
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
		fmt.Println(err)
		return 0
	}
	return f
}

func Verify(rate []string) (bool, error) {
	var err []string
	isValid := true
	for i, value := range rate {
		if i > 2 {
			if isEmpty(value) {
				isValid = false
				err = append(err, fmt.Sprintf("Failed to verify rate, missing fields. \n rate: %s \n null value on: %d", rate, i))
			}
		}
	}
	if !isValid {
		return false, errors.New(fmt.Sprint(err))
	}
	return isValid, nil
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
		fmt.Println(match)
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
	segs := strings.Split(fileName, ".")
	return segs[0]
}
