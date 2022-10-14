package model

import "github.com/shopspring/decimal"

type Wallet struct {
	WalletID string          `json:"wallet_id"`
	Balance  decimal.Decimal `json:"balance"`
}

func (Wallet) TableName() string {
	return "wallets"
}
