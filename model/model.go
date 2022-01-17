package model

import "github.com/shopspring/decimal"

type (
	Wallet struct {
		ID      string          `json:"walletId"`
		Balance decimal.Decimal `json:"balance"`
	}

	Deposit struct {
		ID     string          `json:"walletId"`
		Amount decimal.Decimal `json:"amount"`
	}

	Transfer struct {
		FromID string          `json:"walletFromId"`
		ToID   string          `json:"walletToId"`
		Amount decimal.Decimal `json:"amount"`
	}
)
