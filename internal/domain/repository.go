package domain

import (
	"wallet-api/internal/model"

	"github.com/shopspring/decimal"
)

type Repository interface {
	CreateWallet(wallet *model.Wallet) error

	GetAllWallet() ([]model.Wallet, error)
	GetWallet(walletID string) (*model.Wallet, error)

	DepositWallet(walletID string, amount decimal.Decimal) (*model.Wallet, error)
	TransferWallet(fromWalletID, toWalletID string, amount decimal.Decimal) (*model.Wallet, error)

	DeleteWallet(walletID string) error
}
