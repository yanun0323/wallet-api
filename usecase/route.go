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
	wallet, err := h.repo.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}
	return c.JSON(http.StatusOK, *wallet)
}

func (h *Route) CreateWallet(c echo.Context) error {
	w := &domain.Wallet{}
	if c.Bind(w) != nil {
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
	if c.Bind(d) != nil || d.Amount.IsNegative() {
		return c.JSON(http.StatusBadRequest, nil)
	}

	w, err := h.repo.Get(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}

	w.Balance = w.Balance.Add(d.Amount)
	err2 := h.repo.Update(w)
	if err2 != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return c.JSON(http.StatusOK, *w)
}

func (h *Route) TransferWallet(c echo.Context) error {
	t := new(domain.Transfer)
	if c.Bind(t) != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	wFrom, err := h.repo.Get(t.FromID)
	wTo, err2 := h.repo.Get(t.ToID)
	if err != nil || err2 != nil {
		return c.JSON(http.StatusNotFound, "Can't find wallet.")
	}

	if !isWalletBalanceEnough(wFrom, t.Amount) {
		return c.JSON(http.StatusBadRequest, nil)
	}

	wFrom.Balance = wFrom.Balance.Sub(t.Amount)
	wTo.Balance = wTo.Balance.Add(t.Amount)

	err3 := h.repo.Update(wFrom, wTo)

	if err3 != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, *wFrom)
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
