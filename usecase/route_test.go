package usecase

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-api/domain"

	"github.com/labstack/echo/v4"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func Test_GetAllWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("GetAllWallet", c).Return(nil)
	mockRoute.GetAllWallet(c)
	mockRoute.AssertExpectations(t)
}
func Test_GetWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("GetWallet", c).Return(nil)
	mockRoute.GetWallet(c)
	mockRoute.AssertExpectations(t)
}
func Test_CreateWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("CreateWallet", c).Return(nil)
	mockRoute.CreateWallet(c)
	mockRoute.AssertExpectations(t)
}
func Test_DepositWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("DepositWallet", c).Return(nil)
	mockRoute.DepositWallet(c)
	mockRoute.AssertExpectations(t)
}
func Test_TransferWallet(t *testing.T) {
	mockRoute := new(MockRoute)
	c := createEmptyContext()
	mockRoute.On("TransferWallet", c).Return(nil)
	mockRoute.TransferWallet(c)
	mockRoute.AssertExpectations(t)
}

func Test_DeleteWallet(t *testing.T) {
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
