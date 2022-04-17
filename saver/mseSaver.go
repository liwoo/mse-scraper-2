package saver

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type MSESaver struct {
	FileUrl   string
	ErrorPath string
	Db        *pgxpool.Pool
}

func (u MSESaver) Perform() {
	rows, err := csvToArray(u.FileUrl)

	if err != nil {
		log.Fatal(err)
	}

	copyCount, err := u.Db.CopyFrom(
		context.Background(),
		pgx.Identifier{"people"},
		[]string{"first_name", "last_name", "age"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Finished copying", copyCount, u.FileUrl)
}

func csvToArray(filename string) ([][]interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	rows := [][]interface{}{}

	for _, word := range records {
		var iter []interface{}
		iter = append(iter,
			word[0],
			word[1],
			word[2],
			word[3],
			parseFloat(word[4]),
			parseFloat(word[5]),
			parseFloat(word[6]),
			parseFloat(word[7]),
			parseFloat(word[8]),
			parseFloat(word[9]),
			parseFloat(word[10]),
			parseFloat(word[11]),
			parseFloat(word[12]),
			parseFloat(word[13]),
			parseFloat(word[14]),
			parseFloat(word[16]),
		)
		rows = append(rows, iter)
	}

	return rows, nil
}

func parseFloat(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return f
}
