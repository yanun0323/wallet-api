package main

import (
	"wallet-api/database"
	"wallet-api/handler"
	"wallet-api/model"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	for {
		dsn := "Yanun:Yanun840323@tcp(database.c6ocv0719zbq.ap-northeast-1.rds.amazonaws.com:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			database.Init(db)
			break
		}
	}
}

func main() {
	e := echo.New()

	database.Db.AutoMigrate(&model.Wallet{})
	h := handler.NewHandler(database.Db)

	e.GET("/wallet", h.GetAllWallet)
	e.POST("/wallet/", h.CreateWallet)
	e.PUT("/wallet/", h.TransferWallet)
	e.GET("/wallet/:walletId", h.GetWallet)
	e.PUT("/wallet/:walletId", h.DepositWallet)
	e.DELETE("/wallet/:walletId", h.DeleteWallet)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
