package model

import "github.com/shopspring/decimal"

type Wallet struct {
	WalletID string          `json:"wallet_id" gorm:"unique;not null;index"`
	Balance  decimal.Decimal `json:"balance" gorm:"not null"`
}

func (Wallet) TableName() string {
	return "wallets"
}
