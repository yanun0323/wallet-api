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
	DeleteWallet(c echo.Context) error
}

func NewHandler(db *gorm.DB) Handler {
	var h Handler
	h.db = db
	return h
}

type Handler struct {
	db *gorm.DB
}

func (h *Handler) GetAllWallet(c echo.Context) error {
	result := []model.Wallet{}
	h.db.Table("wallets").Find(&result)
	if result == nil {
		return c.JSON(http.StatusNotFound, nil)
	}
	return c.JSON(http.StatusOK, result)
}

func (h *Handler) GetWallet(c echo.Context) error {
	id := c.Param("walletId")
	w := &model.Wallet{}
	if h.db.First(w, id).Error != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}
	return c.JSON(http.StatusOK, *w)
}

func (h *Handler) CreateWallet(c echo.Context) error {
	u := &model.Wallet{}
	w := &model.Wallet{}
	if c.Bind(u) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if h.db.First(w, u.ID).Error == nil {
		return c.JSON(http.StatusConflict, "Wallet does already exist.")
	}

	if h.db.Create(u).Error != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, *w)
}

func (h *Handler) DepositWallet(c echo.Context) error {
	id := c.Param("walletId")
	m := new(model.Deposit)
	if c.Bind(m) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	w := &model.Wallet{}
	if h.db.First(w, id).Error != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}

	if m.Amount.IsNegative() {
		return c.JSON(http.StatusBadRequest, nil)
	}

	w.Balance = w.Balance.Add(m.Amount)
	h.db.Save(w)
	return c.JSON(http.StatusOK, *w)
}

func (h *Handler) TransferWallet(c echo.Context) error {
	u := new(model.Transfer)
	if c.Bind(u) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	wFrom := &model.Wallet{}
	wTo := &model.Wallet{}

	if !isIdExist(u.FromID, wFrom, h.db) || !isIdExist(u.ToID, wTo, h.db) {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}

	if !isWalletBalanceEnough(wFrom, u.Amount) {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
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

	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, *wFrom)
}

func (h *Handler) DeleteWallet(c echo.Context) error {
	id := c.Param("walletId")
	w := &model.Wallet{}
	if h.db.First(w, id).Error != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}
	if h.db.Delete(w).Error != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func isIdExist(id string, w *model.Wallet, db *gorm.DB) bool {
	return db.First(w, id).Error != nil
}

func isWalletBalanceEnough(w *model.Wallet, amount decimal.Decimal) bool {
	return amount.GreaterThanOrEqual(decimal.NewFromInt32(0)) && w.Balance.GreaterThanOrEqual(amount)
}
