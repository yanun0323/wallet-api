package repository

import (
	"wallet-api/domain"

	"github.com/shopspring/decimal"
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

func (m *MysqlDB) Create(w *domain.Wallet) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		//got, err := GetForUpdate(tx, w.ID)
		got := &domain.Wallet{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(got, w.ID).Error
		if err == nil || got.ID == w.ID {
			return gorm.ErrRegistered
		}
		err = tx.Create(w).Error
		if err != nil {
			return err
		}
		return nil

	})
}

func (m *MysqlDB) Deposit(id string, amount decimal.Decimal) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		//got, err := GetForUpdate(tx, id)
		got := &domain.Wallet{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(got, id).Error
		if err != nil {
			return err
		}
		got.Balance = got.Balance.Add(amount)
		err = tx.Save(got).Error
		if err != nil {
			return err
		}
		return nil
	})
}

func (m *MysqlDB) Transfer(t *domain.Transfer) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		//from, err := GetForUpdate(tx, t.FromID)
		from := &domain.Wallet{}
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(from, t.FromID).Error
		if err != nil {
			return err
		}
		//to, err := GetForUpdate(tx, t.ToID)
		to := &domain.Wallet{}
		err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(to, t.ToID).Error
		if err != nil {
			return err
		}
		if from.Balance.LessThan(t.Amount) {
			return gorm.ErrInvalidData
		}
		from.Balance = from.Balance.Sub(t.Amount)
		to.Balance = to.Balance.Add(t.Amount)
		err = tx.Save(from).Error
		if err != nil {
			return err
		}
		err = tx.Save(to).Error
		if err != nil {
			return err
		}
		return nil
	})
}

// func (m *MysqlDB) Update(ws ...*domain.Wallet) error {
// 	return m.db.Transaction(func(tx *gorm.DB) error {
// 		for _, w := range ws {
// 			err := tx.Save(w).Error
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return nil
// 	})
// }

func (m *MysqlDB) Delete(id string) error {
	w := &domain.Wallet{}
	if err := m.db.First(w, id).Error; err != nil {
		return err
	}
	err := m.db.Delete(w).Error
	return err
}
