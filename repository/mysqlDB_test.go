package repository

import (
	"testing"
	"wallet-api/domain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
)

type MockMysqlDB struct {
	mock.Mock
}

func (m *MockMysqlDB) GetAll() (*[]domain.Wallet, error) {
	args := m.Called()
	return args.Get(0).(*[]domain.Wallet), args.Error(1)
}
func (m *MockMysqlDB) Get(id string) (*domain.Wallet, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Wallet), args.Error(1)
}

func (m *MockMysqlDB) Create(w *domain.Wallet) error {
	args := m.Called(w)
	return args.Error(0)
}
func (m *MockMysqlDB) Deposit(id string, amount decimal.Decimal) error {
	args := m.Called(id, amount)
	return args.Error(0)
}
func (m *MockMysqlDB) Transfer(t *domain.Transfer) error {
	args := m.Called(t)
	return args.Error(0)
}
func (m *MockMysqlDB) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func Test_GetAll(t *testing.T) {
	mysqlDB := new(MockMysqlDB)
	mysqlDB.On("GetAll").Return(&[]domain.Wallet{}, nil)
	mysqlDB.GetAll()

	mysqlDB.AssertExpectations(t)
}
func Test_Get(t *testing.T) {
	mysqlDB := new(MockMysqlDB)
	mysqlDB.On("Get", "123456789").Return(&domain.Wallet{}, nil)
	mysqlDB.Get("123456789")

	mysqlDB.AssertExpectations(t)
}
func Test_Create(t *testing.T) {
	mysqlDB := new(MockMysqlDB)

	w := &domain.Wallet{
		ID:      "123456789",
		Balance: decimal.New(100, 0),
	}

	mysqlDB.On("Create", w).Return(nil)
	mysqlDB.Create(w)

	mysqlDB.AssertExpectations(t)
}

func Test_Deposit(t *testing.T) {
	mysqlDB := new(MockMysqlDB)

	id := "123456789"
	amount := decimal.NewFromInt32(100)
	mysqlDB.On("Deposit", id, amount).Return(nil)
	mysqlDB.Deposit(id, amount)

	mysqlDB.AssertExpectations(t)
}

func Test_Update_TwoWallets(t *testing.T) {
	mysqlDB := new(MockMysqlDB)

	tr := &domain.Transfer{
		FromID: "123456789",
		ToID:   "987654321",
		Amount: decimal.NewFromInt32(100),
	}
	mysqlDB.On("Transfer", tr).Return(nil)
	mysqlDB.Transfer(tr)

	mysqlDB.AssertExpectations(t)
}
func Test_Delete(t *testing.T) {
	mysqlDB := new(MockMysqlDB)
	mysqlDB.On("Delete", "123456789").Return(nil)
	mysqlDB.Delete("123456789")

	mysqlDB.AssertExpectations(t)
}
