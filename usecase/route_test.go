package usecase

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-api/domain"
	"wallet-api/repository"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	w1 = &domain.Wallet{
		ID:      "123456789",
		Balance: decimal.New(100, 0),
	}
	w2 = &domain.Wallet{
		ID:      "987654321",
		Balance: decimal.New(100, 0),
	}
)

func provideDBIRoute(t *testing.T) domain.IRoute {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	require.Nil(t, err)

	resetDBTable(t, db)

	err = db.Create(w1).Error
	require.Nil(t, err)

	err = db.Create(w2).Error
	require.Nil(t, err)

	r := repository.NewMysql(db)
	return NewRoute(r)
}
func provideEmptyDBIRoute(t *testing.T) domain.IRoute {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	require.Nil(t, err)

	resetDBTable(t, db)

	r := repository.NewMysql(db)
	return NewRoute(r)
}

func resetDBTable(t *testing.T, db *gorm.DB) {
	var err error
	err = db.Migrator().DropTable(&domain.Wallet{})
	require.Nil(t, err)
	err = db.AutoMigrate(&domain.Wallet{})
	require.Nil(t, err)
}

func Test_GetAllWallet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet")

	expected := "[]\n"
	route := provideEmptyDBIRoute(t)
	if assert.NoError(t, route.GetAllWallet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet")
	expected = "[{\"walletId\":\"123456789\",\"balance\":\"100\"},{\"walletId\":\"987654321\",\"balance\":\"100\"}]\n"
	route = provideDBIRoute(t)
	if assert.NoError(t, route.GetAllWallet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func Test_GetWallet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//get empty id
	route := provideDBIRoute(t)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("000000000")

	if assert.NoError(t, route.GetWallet(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	//get correct id
	expected := "{\"walletId\":\"123456789\",\"balance\":\"100\"}\n"

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("123456789")

	if assert.NoError(t, route.GetWallet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}
func Test_CreateWallet(t *testing.T) {
	//bad Request
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("123456789")

	route := provideDBIRoute(t)

	if assert.NoError(t, route.CreateWallet(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	//good Request
	requestJSON := "{\"walletId\":\"000000000\",\"balance\":\"100\"}\n"
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//create wallet succeed
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("000000000")

	expected := "{\"walletId\":\"000000000\",\"balance\":\"100\"}\n"

	if assert.NoError(t, route.CreateWallet(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}

	//create existed wallet
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("123456789")

	if assert.NoError(t, route.CreateWallet(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}
func Test_DepositWallet(t *testing.T) {
	//bad Request
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("123456789")

	route := provideDBIRoute(t)

	if assert.NoError(t, route.DepositWallet(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	//good Request
	requestJSON := `{"amount":"100"}`
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//deposit empty wallet
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("000000000")

	if assert.NoError(t, route.DepositWallet(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	//deposit wallet succeed
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("123456789")

	expected := "{\"walletId\":\"123456789\",\"balance\":\"200\"}\n"

	if assert.NoError(t, route.DepositWallet(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}
func Test_TransferWallet(t *testing.T) {
	//bad Request
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet/")

	route := provideDBIRoute(t)

	if assert.NoError(t, route.DepositWallet(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}

	//good Request

	//transfer empty wallet
	requestJSON := `{"walletFromId":"123456789","walletToId": "00000000","amount": 100}`
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet")

	if assert.NoError(t, route.TransferWallet(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	//transfer wallet succeed
	requestJSON = `{"walletFromId":"123456789","walletToId": "987654321","amount": 100}`
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet")

	expected := "{\"walletId\":\"123456789\",\"balance\":\"0\"}\n"

	if assert.NoError(t, route.TransferWallet(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func Test_DeleteWallet(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	//delete empty id
	route := provideDBIRoute(t)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("000000000")

	if assert.NoError(t, route.DeleteWallet(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	//delete correct id

	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	c.SetPath("/wallet/:walletId")
	c.SetParamNames("walletId")
	c.SetParamValues("123456789")

	if assert.NoError(t, route.DeleteWallet(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}

type MockRoute struct {
	mock.Mock
}

func (m *MockRoute) GetAllWallet(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
func (m *MockRoute) CreateWallet(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
func (m *MockRoute) GetWallet(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
func (m *MockRoute) DepositWallet(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
func (m *MockRoute) TransferWallet(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}
func (m *MockRoute) DeleteWallet(c echo.Context) error {
	args := m.Called(c)
	return args.Error(0)
}

func createEmptyContext() echo.Context {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return c
}

func MockTest_GetAllWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("GetAllWallet", c).Return(nil)
	mockRoute.GetAllWallet(c)
	mockRoute.AssertExpectations(t)
}
func MockTest_GetWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("GetWallet", c).Return(nil)
	mockRoute.GetWallet(c)
	mockRoute.AssertExpectations(t)
}
func MockTest_CreateWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("CreateWallet", c).Return(nil)
	mockRoute.CreateWallet(c)
	mockRoute.AssertExpectations(t)
}
func MockTest_DepositWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("DepositWallet", c).Return(nil)
	mockRoute.DepositWallet(c)
	mockRoute.AssertExpectations(t)
}
func MockTest_TransferWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("TransferWallet", c).Return(nil)
	mockRoute.TransferWallet(c)
	mockRoute.AssertExpectations(t)
}

func MockTest_DeleteWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("DeleteWallet", c).Return(nil)
	mockRoute.DeleteWallet(c)
	mockRoute.AssertExpectations(t)
}

func Test_isWalletBalanceEnough(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		input1   *domain.Wallet
		input2   decimal.Decimal
		expected bool
	}{
		{input1: &domain.Wallet{ID: "123456789", Balance: decimal.NewFromInt32(0)}, input2: decimal.NewFromInt32(0), expected: true},
		{input1: &domain.Wallet{ID: "123456789", Balance: decimal.NewFromInt32(100)}, input2: decimal.NewFromInt32(100), expected: true},
		{input1: &domain.Wallet{ID: "123456789", Balance: decimal.NewFromInt32(100)}, input2: decimal.NewFromInt32(50), expected: true},
		{input1: &domain.Wallet{ID: "123456789", Balance: decimal.NewFromInt32(0)}, input2: decimal.NewFromInt32(50), expected: false},
		{input1: &domain.Wallet{ID: "123456789", Balance: decimal.NewFromInt32(50)}, input2: decimal.NewFromInt32(100), expected: false},
	}

	for _, test := range tests {
		assert.Equal(isWalletBalanceEnough(test.input1, test.input2), test.expected)
	}
}
