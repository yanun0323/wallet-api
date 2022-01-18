package main

import (
	"wallet-api/delivery"
	"wallet-api/domain"
	"wallet-api/repository"
	"wallet-api/usecase"

	_ "github.com/joho/godotenv/autoload"
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

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
