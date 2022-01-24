package usecase

import (
	"net/http"
	"wallet-api/domain"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
)

type Route struct {
	repo domain.IRepository
}

func NewRoute(r domain.IRepository) domain.IRoute {
	return &Route{
		repo: r,
	}
}

func (h *Route) GetAllWallet(c echo.Context) error {
	wallet, err := h.repo.GetAll()
	if err != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}
	return c.JSON(http.StatusOK, wallet)
}

func (h *Route) GetWallet(c echo.Context) error {
	id := c.Param("walletId")
	if id == "" {
		return c.JSON(http.StatusNotFound, "WalletId can't be empty.")
	}
	wallet, err := h.repo.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}
	return c.JSON(http.StatusOK, *wallet)
}

func (h *Route) CreateWallet(c echo.Context) error {
	w := &domain.Wallet{}
	if c.Bind(w) != nil || w.ID == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}
	err := h.repo.Create(w)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, *w)
}

func (h *Route) DepositWallet(c echo.Context) error {
	id := c.Param("walletId")
	d := new(domain.Deposit)
	if c.Bind(d) != nil || !d.Amount.IsPositive() {
		return c.JSON(http.StatusBadRequest, nil)
	}

	if h.repo.Deposit(id, d.Amount) != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *Route) TransferWallet(c echo.Context) error {
	t := new(domain.Transfer)
	if c.Bind(t) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}
	if !t.Amount.IsPositive() {
		return c.JSON(http.StatusBadRequest, nil)
	}
	err := h.repo.Transfer(t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusCreated, nil)
}

func (h *Route) DeleteWallet(c echo.Context) error {
	id := c.Param("walletId")
	err := h.repo.Delete(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}

	return c.JSON(http.StatusNoContent, nil)
}

func isWalletBalanceEnough(w *domain.Wallet, amount decimal.Decimal) bool {
	return amount.GreaterThanOrEqual(decimal.NewFromInt32(0)) && w.Balance.GreaterThanOrEqual(amount)
}
