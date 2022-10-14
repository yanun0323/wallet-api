package service

import (
	"net/http"
	"wallet-api/internal/model"
	"wallet-api/internal/util"

	"github.com/labstack/echo/v4"
)

func (svc *Service) GetAllWallet(c echo.Context) error {
	wallets, err := svc.repo.GetAllWallet()
	if err != nil {
		svc.l.Errorf("get all wallet error, %+v", err)
		return c.JSON(http.StatusNotFound, util.Msg("can't find wallet."))
	}
	svc.l.Infof("successfully get all wallet")
	return c.JSON(http.StatusOK, wallets)
}

func (svc *Service) GetWallet(c echo.Context) error {
	id := c.Param("walletID")
	if len(id) == 0 {
		svc.l.Warn("get wallet error, empty wallet ID")
		return c.JSON(http.StatusNotFound, util.Msg("empty wallet ID."))
	}
	wallet, err := svc.repo.GetWallet(id)
	if err != nil {
		svc.l.Errorf("get wallet error, %+v", err)
		return c.JSON(http.StatusNotFound, util.Msg("can't find wallet."))
	}

	svc.l.Infof("successfully get wallet: %s", id)
	return c.JSON(http.StatusOK, *wallet)
}

func (svc *Service) CreateWallet(c echo.Context) error {
	w := &model.Wallet{}
	if c.Bind(w) != nil || len(w.WalletID) == 0 {
		svc.l.Warn("create wallet error, bad request")
		return c.JSON(http.StatusBadRequest, util.Msg("bad request."))
	}
	err := svc.repo.CreateWallet(w)
	if err != nil {
		svc.l.Errorf("create wallet error, %+v", err)
		return c.JSON(echo.ErrInternalServerError.Code, util.Msg("unknown error."))
	}

	svc.l.Infof("successfully create wallet: %s", w.WalletID)
	return c.JSON(http.StatusOK, w)
}

func (svc *Service) DepositWallet(c echo.Context) error {
	id := c.Param("walletId")
	if len(id) == 0 {
		svc.l.Warn("deposit wallet error, empty wallet ID")
		return c.JSON(http.StatusNotFound, util.Msg("empty wallet ID."))
	}

	d := &model.Deposit{}
	if c.Bind(d) != nil || !d.Amount.IsPositive() {
		svc.l.Warn("deposit wallet error, bad request")
		return c.JSON(http.StatusBadRequest, util.Msg("bad request."))
	}
	wallet, err := svc.repo.DepositWallet(id, d.Amount)
	if err != nil {
		svc.l.Errorf("deposit wallet error, %+v", err)
		return c.JSON(http.StatusInternalServerError, util.Msg("unknown error."))
	}

	svc.l.Infof("successfully deposit wallet: %s", id)
	return c.JSON(http.StatusOK, *wallet)
}

func (svc *Service) TransferWallet(c echo.Context) error {
	t := &model.Transfer{}
	if c.Bind(t) != nil {
		svc.l.Warn("transfer wallet error, bad request")
		return c.JSON(http.StatusBadRequest, util.Msg("bad request."))
	}
	if !t.Amount.IsPositive() {
		svc.l.Warn("transfer amount is less than zero")
		return c.JSON(http.StatusBadRequest, util.Msg("amount must be greater than zero."))
	}
	fromWallet, err := svc.repo.TransferWallet(t.FromWalletID, t.ToWalletID, t.Amount)
	if err != nil {
		svc.l.Errorf("transfer wallet error, %+v", err)
		return c.JSON(http.StatusInternalServerError, util.Msg("unknown error."))
	}

	svc.l.Infof("successfully transfer wallet %s to %s", t.FromWalletID, t.ToWalletID)
	return c.JSON(http.StatusOK, *fromWallet)
}

func (svc *Service) DeleteWallet(c echo.Context) error {
	id := c.Param("walletId")
	if len(id) == 0 {
		svc.l.Warn("delete wallet error, empty wallet ID")
		return c.JSON(http.StatusNotFound, util.Msg("empty wallet ID."))
	}
	err := svc.repo.DeleteWallet(id)
	if err != nil {
		svc.l.Errorf("delete wallet error, %+v", err)
		return c.JSON(http.StatusNotFound, util.Msg("can't find wallet."))
	}

	svc.l.Infof("successfully delete wallet: %s", id)
	return c.JSON(http.StatusOK, util.Msg("success."))
}
