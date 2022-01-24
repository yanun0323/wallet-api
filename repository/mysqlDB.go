package repository

import (
	"wallet-api/domain"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MysqlDB struct {
	db *gorm.DB
}

func NewMysql(db *gorm.DB) domain.IRepository {
	return &MysqlDB{db}
}

func (m *MysqlDB) GetAll() (*[]domain.Wallet, error) {
	result := []domain.Wallet{}
	err := m.db.Table("wallets").Find(&result).Error
	return &result, err
}

func (m *MysqlDB) Get(id string) (*domain.Wallet, error) {
	w := &domain.Wallet{}
	err := m.db.First(w, id).Error
	return w, err
}

//.Clauses(clause.Locking{Strength: "UPDATE"})
func (m *MysqlDB) GetForUpdate(id string) (*domain.Wallet, error) {
	w := &domain.Wallet{}
	err := m.db.Clauses(clause.Locking{Strength: "UPDATE"}).First(w, id).Error
	return w, err
}

func (m *MysqlDB) Create(w *domain.Wallet) error {
	err := m.db.Create(w).Error
	return err
}

func (m *MysqlDB) Update(ws ...*domain.Wallet) error {

	return m.db.Transaction(func(tx *gorm.DB) error {
		for _, w := range ws {
			err := tx.Save(w).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *MysqlDB) Delete(id string) error {
	w := &domain.Wallet{}
	if err := m.db.First(w, id).Error; err != nil {
		return err
	}
	err := m.db.Delete(w).Error
	return err
}
