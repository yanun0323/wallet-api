package repository

import (
	"testing"
	"wallet-api/domain"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB     *gorm.DB
	wallet *domain.Wallet
}

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
func (m *MockMysqlDB) Update(ws ...*domain.Wallet) error {
	args := m.Called(ws)
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

func Test_Update(t *testing.T) {
	mysqlDB := new(MockMysqlDB)

	w := &domain.Wallet{
		ID:      "123456789",
		Balance: decimal.NewFromInt32(100),
	}
	expected := []*domain.Wallet{w}
	mysqlDB.On("Update", expected).Return(nil)
	mysqlDB.Update(w)

	mysqlDB.AssertExpectations(t)
}

func Test_Update_TwoWallets(t *testing.T) {
	mysqlDB := new(MockMysqlDB)

	w1 := &domain.Wallet{
		ID:      "123456789",
		Balance: decimal.NewFromInt32(100),
	}
	w2 := &domain.Wallet{
		ID:      "987654321",
		Balance: decimal.New(100, 0),
	}
	expected := []*domain.Wallet{w1, w2}
	mysqlDB.On("Update", expected).Return(nil)
	mysqlDB.Update(w1, w2)

	mysqlDB.AssertExpectations(t)
}
func Test_Delete(t *testing.T) {
	mysqlDB := new(MockMysqlDB)
	mysqlDB.On("Delete", "123456789").Return(nil)
	mysqlDB.Delete("123456789")

	mysqlDB.AssertExpectations(t)
}
