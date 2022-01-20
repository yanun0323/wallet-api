package main

import (
	"net/http"
	"wallet-api/delivery"
	"wallet-api/domain"
	"wallet-api/repository"
	"wallet-api/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	_db, err := gorm.Open(sqlite.Open("wallets.db"), &gorm.Config{})
	db = _db
	if err != nil {
		panic("failed to connect database")
	}
}

func main() {
	e := echo.New()
	db.AutoMigrate(&domain.Wallet{})

	repo := repository.NewMysql(db)
	usecase := usecase.NewRoute(repo)
	delivery.NewHandler(e, usecase)

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "Hello Bing Chilling.")
	})

	// Start server
	e.Logger.Fatal(e.Start(":8080"))

}
