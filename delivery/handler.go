package delivery

import (
	model "wallet-api/domain"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	usecase model.IRoute
}

func NewHandler(e *echo.Echo, u model.IRoute) {
	handler := &Handler{
		usecase: u,
	}
	e.GET("/wallet", handler.usecase.GetAllWallet)
	e.POST("/wallet/", handler.usecase.CreateWallet)
	e.PUT("/wallet/", handler.usecase.TransferWallet)
	e.GET("/wallet/:walletId", handler.usecase.GetWallet)
	e.PUT("/wallet/:walletId", handler.usecase.DepositWallet)
	e.DELETE("/wallet/:walletId", handler.usecase.DeleteWallet)
}
