package clean

import (
	"context"
	"fmt"
	"log"
	"mseScraping/data"
	"mseScraping/pkg/conf"
	"net/http"
	"time"

	"database/sql"
	"github.com/sirsean/go-pool"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func CleanDatabase(w http.ResponseWriter, r *http.Request, conf conf.Conf) {
	cleanDb(conf)
	w.Write([]byte("Clean Db"))
}

func cleanDb(s conf.Conf) {
	fmt.Println("Saving to Database from... ", s.CleanedJsonPath)
	p := pool.NewPool(s.QueueSize, s.WorkerNum)

	var pgconn *pgdriver.Connector = pgdriver.NewConnector(pgdriver.WithDSN(s.DBConnectionString))
	psdb := sql.OpenDB(pgconn)
	db := bun.NewDB(psdb, pgdialect.New())

	p.Start()
	cleanDatabase(db)
	p.Close()
}

func cleanDatabase(db *bun.DB) {
	layout := "2006-01-02"
	rates := []data.DailyCompanyRateModel{}

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
		previousRate := new(data.DailyCompanyRateModel)
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
