package domain

import (
	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type (
	Wallet struct {
		ID      string          `json:"walletId"`
		Balance decimal.Decimal `json:"balance"`
	}

	Deposit struct {
		Amount decimal.Decimal `json:"amount"`
	}

	Transfer struct {
		FromID string          `json:"walletFromId"`
		ToID   string          `json:"walletToId"`
		Amount decimal.Decimal `json:"amount"`
	}
)

type IRepository interface {
	GetAll() (*[]Wallet, error)
	Get(id string) (*Wallet, error)
	Create(w *Wallet) error
	Deposit(id string, amount decimal.Decimal) error
	Transfer(t *Transfer) error
	//Update(w ...*Wallet) error
	Delete(id string) error
}

type IRoute interface {
	GetAllWallet(c echo.Context) error
	CreateWallet(c echo.Context) error
	GetWallet(c echo.Context) error
	DepositWallet(c echo.Context) error
	TransferWallet(c echo.Context) error
	DeleteWallet(c echo.Context) error
}
