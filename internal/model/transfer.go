package model

import "github.com/shopspring/decimal"

type Transfer struct {
	FromWalletID string          `json:"from_wallet_id"`
	ToWalletID   string          `json:"to_wallet_id"`
	Amount       decimal.Decimal `json:"amount"`
}
