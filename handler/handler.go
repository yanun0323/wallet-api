package handler

import (
	"net/http"
	"wallet-api/model"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type IHandler interface {
	GetAllWallet(c echo.Context) error
	CreateWallet(c echo.Context) error
	GetWallet(c echo.Context) error
	DepositWallet(c echo.Context) error
	TransferWallet(c echo.Context) error
}

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) Handler {
	var h Handler
	h.db = db
	return h
}

func (h *Handler) GetAllWallet(c echo.Context) error {
	result := []model.Wallet{}
	h.db.Table("wallets").Find(&result)
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) GetWallet(c echo.Context) error {
	id := c.Param("walletId")
	w := &model.Wallet{}
	if h.db.First(w, id).Error != nil {
		return c.JSON(http.StatusOK, "user does not exit!")
	}
	return c.JSON(http.StatusOK, *w)
}

func (h *Handler) CreateWallet(c echo.Context) error {
	u := &model.Wallet{}

	if err := c.Bind(u); err != nil {
		return err
	}

	if h.db.First(&model.Wallet{}, u.ID).Error == nil {
		return c.JSON(http.StatusOK, "user already exit!")
	}

	if err := h.db.Create(u).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, "Succeesed")
}

func (h *Handler) DepositWallet(c echo.Context) error {
	u := new(model.Deposit)
	if err := c.Bind(u); err != nil {
		return err
	}
	w := &model.Wallet{}
	if h.db.First(w, u.ID).Error != nil {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	if u.Amount.IsNegative() {
		return c.JSON(http.StatusOK, "amount can't be negative!")
	}

	w.Balance = w.Balance.Add(u.Amount)
	h.db.Save(w)
	return c.JSON(http.StatusOK, *w)
}

func (h *Handler) TransferWallet(c echo.Context) error {
	u := new(model.Transfer)
	if err := c.Bind(u); err != nil {
		return err
	}

	wFrom := &model.Wallet{}
	wTo := &model.Wallet{}

	if !isIdExist(u.FromID, wFrom, h.db) || !isIdExist(u.ToID, wTo, h.db) {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	if !isWalletBalanceEnough(wFrom, u.Amount) {
		return c.JSON(http.StatusOK, "user balance is not enough!")
	}

	h.db.Transaction(func(tx *gorm.DB) error {
		wFrom.Balance = wFrom.Balance.Sub(u.Amount)
		wTo.Balance = wTo.Balance.Add(u.Amount)

		if err := h.db.Save(wFrom).Error; err != nil {
			return err
		}
		if err := h.db.Save(wTo).Error; err != nil {
			return err
		}
		return nil
	})

	return c.JSON(http.StatusOK, *wFrom)
}

func isIdExist(id string, w *model.Wallet, db *gorm.DB) bool {
	return db.First(w, id).Error != nil
}

func isWalletBalanceEnough(w *model.Wallet, amount decimal.Decimal) bool {
	return amount.GreaterThanOrEqual(decimal.NewFromInt32(0)) && w.Balance.GreaterThanOrEqual(amount)
}
