package main

import (
	"net/http"
	"wallet-api/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

//----------
// Handlers
//----------

func createWallet(c echo.Context) error {
	u := new(model.Wallet)
	if err := c.Bind(u); err != nil {
		return err
	}

	found := db.First(&model.Wallet{}, u.ID)
	if found.Error == nil {
		return c.JSON(http.StatusOK, "user already exit!")
	}

	if db.Create(&u).Error != nil {
		return c.JSON(http.StatusOK, found.Error)
	}

	return c.JSON(http.StatusOK, "Succeesed")
}

func getWallet(c echo.Context) error {
	id := c.Param("walletId")

	w := &model.Wallet{}
	if db.First(w, id).Error != nil {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	return c.JSONPretty(http.StatusOK, *w, " ")
}

func depositWallet(c echo.Context) error {
	u := new(model.Deposit)
	if err := c.Bind(u); err != nil {
		return err
	}

	w := &model.Wallet{}
	if db.First(w, u.ID).Error != nil {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	if u.Amount.IsNegative() {
		return c.JSON(http.StatusOK, "amount can't be negative!")
	}

	w.Balance = w.Balance.Add(u.Amount)
	db.Save(w)
	return c.JSONPretty(http.StatusOK, *w, " ")
}

func isWalletBalanceEnough(w *model.Wallet, amount decimal.Decimal) bool {
	return amount.GreaterThanOrEqual(decimal.NewFromInt32(0)) && w.Balance.GreaterThanOrEqual(amount)
}

func transferWallet(c echo.Context) error {
	u := new(model.Transfer)
	if err := c.Bind(u); err != nil {
		return err
	}
	wFrom := &model.Wallet{}
	wTo := &model.Wallet{}

	if db.First(wFrom, u.FromID).Error != nil || db.First(wTo, u.ToID).Error != nil {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	if !isWalletBalanceEnough(wFrom, u.Amount) {
		return c.JSON(http.StatusOK, "user balance is not enough!")
	}

	db.Transaction(func(tx *gorm.DB) error {
		wFrom.Balance = wFrom.Balance.Sub(u.Amount)
		wTo.Balance = wTo.Balance.Add(u.Amount)

		if err := db.Save(wFrom).Error; err != nil {
			return err
		}
		if err := db.Save(wTo).Error; err != nil {
			return err
		}
		return nil
	})

	return c.JSONPretty(http.StatusOK, *wFrom, " ")
}

func getAllWallet(c echo.Context) error {
	result := []model.Wallet{}
	db.Table("wallets").Find(&result)
	return c.JSONPretty(http.StatusOK, result, " ")
}

func main() {
	e := echo.New()
	dsn := "Yanun:Yanun840323@tcp(database.c6ocv0719zbq.ap-northeast-1.rds.amazonaws.com:3306)/wallet?charset=utf8mb4&parseTime=True&loc=Local"
	walletDb, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db = walletDb
	//db.Migrator().DropTable(&wallet{})
	db.AutoMigrate(&model.Wallet{})

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/wallet", getAllWallet)
	e.GET("/wallet/:walletId", getWallet)
	e.POST("/wallet/", createWallet)
	e.PUT("/wallet/deposit", depositWallet)
	e.PUT("/wallet/transfer", transferWallet)
	// e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
