package clean

import (
	"fmt"
	"mseScraping/api/clean/dbCleaner"
	"mseScraping/pkg/conf"
	"net/http"

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

	dbCleaner.CleanDatabase(db)

	p.Close()
}
