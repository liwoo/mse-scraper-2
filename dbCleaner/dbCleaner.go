package dbCleaner

import (
	"context"
	"fmt"
	"log"
	"mseScraping/saver"
	"time"

	"github.com/uptrace/bun"
)

func CleanDatabase(db *bun.DB) {
	layout := "2006-01-02"
	rates := []saver.DailyCompanyRateModel{}

	err := db.NewSelect().Model(&rates).Where("buy = ?", -1).WhereOr("sell = ?", -1).Order("date ASC").Scan(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	for _, rate := range rates {
		date, err := time.Parse(layout, rate.DATE)

		if err != nil {
			fmt.Println(err)
			continue
		}

		previousDate := date.AddDate(0, 0, -1)
		previousRate := new(saver.DailyCompanyRateModel)
		colErr := db.NewSelect().Model(previousRate).Where("date = ?", previousDate.Format(layout)).Where("code = ?", rate.CODE).Scan(context.Background())

		if colErr != nil {
			fmt.Println(colErr)
			continue
		}

		res, insErr := db.NewUpdate().Model(&rate).
			Set("buy = ?", previousRate.BUY).
			Set("sell = ?", previousRate.SELL).
			WherePK().Exec(context.Background())
		if insErr != nil {
			fmt.Println(insErr)
		}

		fmt.Println(res.RowsAffected())
	}

}
