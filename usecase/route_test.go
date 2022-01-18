package usecase

import (
	"testing"
	"wallet-api/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewRoute(t *testing.T) {
	n := new(domain.IRepository)
	r := NewRoute(*n)
	assert.NotNil(t, r)
	assert.NotNil(t, r.CreateWallet)
}
