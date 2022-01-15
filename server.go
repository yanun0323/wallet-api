package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	user struct {
		ID      string `json:"id"`
		Balance int    `json:"balance"`
		// Biils   []bill `json:"bills"`
	}

	operation struct {
		Amount int `json:"amount"`
	}

	// bill struct {
	// 	Date    string `json:"date"`
	// 	Reciver string `json:"reciver"`
	// 	Amount  int    `json:"amount"`
	// 	Sender  string `json:"sender"`
	// }
)

var (
	users = map[string]*user{}
)

//----------
// Handlers
//----------

func createUser(c echo.Context) error {
	id := c.Param("id")

	_, exist := users[id]
	if exist {
		return c.JSON(http.StatusOK, "user already exit!")
	}

	u := &user{
		ID:      id,
		Balance: 0,
	}

	users[u.ID] = u
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	data, exist := users[id]
	if !exist {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	return c.JSONPretty(http.StatusOK, data, " ")
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	u := new(operation)
	if err := c.Bind(u); err != nil {
		return err
	}

	data, exist := users[id]
	if !exist {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	if !isUserBalanceEnough(id, u.Amount) {
		return c.JSON(http.StatusOK, "user balance is not enough!")
	}

	addUserBalance(id, u.Amount)
	// addUserBill(id, CreateBill("-", id, u.Amount))
	return c.JSONPretty(http.StatusOK, data, " ")
}

func isUserBalanceEnough(id string, amount int) bool {
	return amount >= 0 || users[id].Balance >= -amount
}

func addUserBalance(id string, amount int) {
	users[id].Balance += amount
}

// func addUserBill(id string, b bill) {
// 	users[id].Biils = append(users[id].Biils, b)
// }
// func CreateBill(sender string, resiver string, amount int) bill {
// 	return bill{Sender: sender, Reciver: resiver, Amount: amount, Date: time.Now().UTC().GoString()}
// }

func transferMoney(c echo.Context) error {
	sender := c.Param("id")
	resiver := c.Param("resiver")
	u := new(operation)
	if err := c.Bind(u); err != nil {
		return err
	}
	data, senderExist := users[sender]
	_, resiverExist := users[resiver]

	if u.Amount > 0 {
		return c.JSON(http.StatusOK, "amount can't be greater than zero!")
	}

	if !senderExist || !resiverExist {
		return c.JSON(http.StatusOK, "user does not exit!")
	}

	if !isUserBalanceEnough(sender, u.Amount) {
		return c.JSON(http.StatusOK, "user balance is not enough!")
	}
	addUserBalance(sender, u.Amount)
	addUserBalance(resiver, -u.Amount)
	// bill := CreateBill(sender, resiver, u.Amount)
	// addUserBill(sender, bill)
	// addUserBill(resiver, bill)
	return c.JSONPretty(http.StatusOK, data, " ")
}

// func deleteUser(c echo.Context) error {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	delete(users, id)
// 	return c.NoContent(http.StatusNoContent)
// }

func getAllUsers(c echo.Context) error {
	return c.JSONPretty(http.StatusOK, users, " ")
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", getAllUsers)
	e.POST("/users/:id", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.PUT("/users/:id/:resiver", transferMoney)
	// e.DELETE("/users/:id", deleteUser)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))

}
