package sqlite

import (
	"fmt"
	"wallet-api/internal/model"

	"github.com/pkg/errors"

	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SqliteDao struct {
	db *gorm.DB
}

func NewSqliteDao(db string) (*SqliteDao, error) {
	client, err := gorm.Open(sqlite.Open(db))
	if err != nil {
		return nil, err
	}

	dao := &SqliteDao{
		db: client,
	}
	dao.Init()

	return dao, nil
}

func (dao *SqliteDao) Init() {
	dao.db.AutoMigrate(&model.Wallet{})
}

func (dao *SqliteDao) CreateWallet(wallet *model.Wallet) error {
	return dao.db.Create(wallet).Error
}

func (dao *SqliteDao) GetAllWallet() ([]model.Wallet, error) {
	ws := []model.Wallet{}
	if err := dao.db.Find(&ws).Error; err != nil {
		return nil, err
	}
	return ws, nil
}

func (dao *SqliteDao) GetWallet(walletID string) (*model.Wallet, error) {
	w := &model.Wallet{}
	err := dao.db.Where("wallet_id=?", walletID).First(w).Error
	if err != nil {
		return nil, err
	}
	return w, err
}

func (dao *SqliteDao) DepositWallet(walletID string, amount decimal.Decimal) (*model.Wallet, error) {
	w := &model.Wallet{}
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		if err := dao.db.Where("wallet_id=?", walletID).First(w).Error; err != nil {
			return err
		}

		w.Balance = w.Balance.Add(amount)
		if err := dao.db.Where("wallet_id=?", w.WalletID).Save(w).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return w, nil
}

func (dao *SqliteDao) TransferWallet(fromWalletID, toWalletID string, amount decimal.Decimal) (*model.Wallet, error) {
	var from, to *model.Wallet
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		wallets := []model.Wallet{}
		err := dao.db.Where("wallet_id IN(?)", []string{fromWalletID, toWalletID}).Find(&wallets).Error
		if err != nil {
			return err
		}

		for i := range wallets {
			switch wallets[i].WalletID {
			case fromWalletID:
				from = &wallets[i]
			case toWalletID:
				to = &wallets[i]
			}
		}

		if from == nil {
			return errors.New(fmt.Sprintf("can't find wallet id, %s", fromWalletID))
		}

		if to == nil {
			return errors.New(fmt.Sprintf("can't find wallet id, %s", toWalletID))
		}

		if from.Balance.LessThan(amount) {
			return errors.New("balance not enough")
		}

		from.Balance.Sub(amount)
		to.Balance.Add(amount)

		if err := dao.db.Where("wallet_id=?", from.WalletID).Save(from).Error; err != nil {
			return err
		}

		if err := dao.db.Where("wallet_id=?", to.WalletID).Save(to).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return from, nil
}

func (dao *SqliteDao) DeleteWallet(walletID string) error {
	return dao.db.Where("wallet_id=?", walletID).Delete(&model.Wallet{}).Error
}
