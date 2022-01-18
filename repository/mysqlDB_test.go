package repository

import (
	"testing"
	"wallet-api/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB     *gorm.DB
	wallet *domain.Wallet
}

type MyMockedObject struct {
	mock.Mock
}

func (m *MyMockedObject) GetAll() (*[]domain.Wallet, error) {
	args := m.Called()

	f := args.Get(0)
	d := f.([]domain.Wallet)

	return &d, args.Error(1)
}

func TestGetAll(t *testing.T) {

	testObj := new(MyMockedObject)
	testObj.On("DoSomething", mock.Anything).Return(true, nil)
	//m := &mysqlDB{}
	testObj.AssertExpectations(t)

	require.NoError(t, nil)
	assert.Nil(t, nil)
}
