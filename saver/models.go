package saver

import "github.com/uptrace/bun"

type DailyCompanyRateModel struct {
	bun.BaseModel `bun:"table:daily_company_rates,alias:u"`

	ID        int64 `bun:",pk,autoincrement"`
	NO        string
	HIGH      string
	LOW       string
	CODE      string
	BUY       float64
	SELL      float64
	PCP       float64
	TCP       float64
	VOL       int64
	DIVNET    float64 `bun:"div_net"`
	DIVYIELD  float64 `bun:"div_yield"`
	EARNYIELD float64 `bun:"earn_yield"`
	PERATIO   float64 `bun:"pe_ratio"`
	PBVRATION float64 `bun:"pbv_ratio"`
	CAP       float64
	PROFIT    float64
	SHARES    float64
	DATE      string
}
