package repository

import (
	"wallet-api/domain"

	"gorm.io/gorm"
)

type mysql struct {
	db *gorm.DB
}

func NewMysql(db *gorm.DB) domain.IRepository {
	return &mysql{db}
}

func (m *mysql) GetAll() (*[]domain.Wallet, error) {
	result := []domain.Wallet{}
	err := m.db.Table("wallets").Find(&result).Error
	return &result, err
}

func (m *mysql) Get(id string) (*domain.Wallet, error) {
	w := &domain.Wallet{}
	err := m.db.First(w, id).Error
	return w, err
}

func (m *mysql) Create(w *domain.Wallet) error {
	err := m.db.Create(w).Error
	return err
}

func (m *mysql) Update(ws ...*domain.Wallet) error {

	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, w := range ws {
			err := m.db.Save(w).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *mysql) Delete(id string) error {
	w := &domain.Wallet{}
	if err := m.db.First(w, id).Error; err != nil {
		return err
	}
	err := m.db.Delete(w).Error
	return err
}
