package app

import (
	"context"
	"net/http"
	"wallet-api/internal/service"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Run() {

	e := echo.New()
	l := e.Logger

	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(10)))
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	svc, err := service.New(context.Background(), l)
	l.Fatal(err)

	SetRouter(e, svc)

	l.Fatal(e.Start(":8080"))

	// sigterm := make(chan os.Signal, 1)
	// signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	// <-sigterm
	// l.Info("shutdown process start")
}

func SetRouter(e *echo.Echo, service *service.Service) {
	e.GET("/wallet", service.GetAllWallet)
	e.POST("/wallet/", service.CreateWallet)
	e.PUT("/wallet/", service.TransferWallet)
	e.GET("/wallet/:walletId", service.GetWallet)
	e.PUT("/wallet/:walletId", service.DepositWallet)
	e.DELETE("/wallet/:walletId", service.DeleteWallet)
}
