package main

import (
	"wallet-api/delivery"
	"wallet-api/domain"
	"wallet-api/repository"
	"wallet-api/usecase"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	e := echo.New()
	db, err := connectDatabase()
	if err != nil {
		e.Logger.Fatal(err)
	}
	db.AutoMigrate(&domain.Wallet{})

	repo := repository.NewMysql(db)
	usecase := usecase.NewRoute(repo)
	delivery.NewHandler(e, usecase)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
func connectDatabase() (*gorm.DB, error) {
	dsn := "Yanun:Yanun840323@tcp(database.c6ocv0719zbq.ap-northeast-1.rds.amazonaws.com:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
