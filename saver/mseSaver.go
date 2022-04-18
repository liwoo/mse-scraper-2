package saver

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/uptrace/bun"
)

type MSESaver struct {
	FileUrl   string
	ErrorPath string
	Db        *bun.DB
}

func (u MSESaver) Perform() {
	rows, codes, err := csvToArray(u.FileUrl)

	if err != nil {
		log.Fatal(err)
	}

	saveCodes(codes, u.Db)

	res, err := u.Db.NewInsert().Model(&rows).On("CONFLICT (id) DO UPDATE").Exec(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		fmt.Println("Error for getting rates affected:", err)
	}

	fmt.Println("Finished copying", affected, u.FileUrl)
}

func saveCodes(codes []string, db *bun.DB) {
	var rows = []CompanyModel{}
	for _, code := range codes {
		rows = append(rows, CompanyModel{
			MSECODE: code,
		})
	}
	res, err := db.NewInsert().Model(&rows).Exec(context.Background())

	if err != nil {
		fmt.Println("Error saving codes:", err)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Finished adding codes", affected)
}

func csvToArray(filePath string) ([]DailyCompanyRateModel, []string, error) {
	fileName := filepath.Base(filePath)
	date := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	f, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	codes := []string{}
	rows := []DailyCompanyRateModel{}

	for _, word := range records {
		var no = word[0]
		var code = word[3]
		var id = strings.ReplaceAll(date, "-", "") + no
		var vers = DailyCompanyRateModel{
			ID:        id,
			NO:        no,
			HIGH:      word[1],
			LOW:       word[2],
			CODE:      code,
			BUY:       parseFloat(word[4]),
			SELL:      parseFloat(word[5]),
			PCP:       parseFloat(word[6]),
			TCP:       parseFloat(word[7]),
			VOL:       parseInt(word[8]),
			DIVNET:    parseFloat(word[9]),
			DIVYIELD:  parseFloat(word[10]),
			EARNYIELD: parseFloat(word[11]),
			PERATIO:   parseFloat(word[12]),
			PBVRATION: parseFloat(word[13]),
			CAP:       parseFloat(word[14]),
			PROFIT:    parseFloat(word[15]),
			SHARES:    parseFloat(word[16]),
			DATE:      date,
		}
		codes = append(codes, code)
		rows = append(rows, vers)
	}

	return rows, codes, nil
}

func parseInt(value string) int64 {
	i, err := strconv.ParseInt(value, 10, 16)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return i
}

func parseFloat(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return f
}
