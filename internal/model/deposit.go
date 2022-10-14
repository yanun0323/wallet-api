package model

import "github.com/shopspring/decimal"

type Deposit struct {
	Amount decimal.Decimal `json:"amount"`
}
