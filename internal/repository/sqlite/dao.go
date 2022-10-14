package sqlite

import (
	"database/sql"
	"fmt"
	"wallet-api/internal/model"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/shopspring/decimal"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	ErrBalanceNotEnough = fmt.Errorf("balance not enough")
	ErrWalletExist      = fmt.Errorf("wallet already exist")
	ErrWalletNotExist   = fmt.Errorf("can't find wallet")
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
	if viper.GetBool("sqlite.debug") {
		dao.db = dao.db.Debug()
	}
	return dao, nil
}

func (dao *SqliteDao) Init() {
	dao.db.AutoMigrate(&model.Wallet{})
}

func (dao *SqliteDao) CreateWallet(wallet *model.Wallet) error {
	return dao.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.Table(wallet.TableName()).Where("wallet_id=?", wallet.WalletID).Count(&count).Error
		if errors.Is(err, sql.ErrNoRows) || count == 0 {
			return tx.Create(wallet).Error
		}

		if err != nil {
			return err
		}

		return ErrWalletExist
	})
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
		if err := tx.Where("wallet_id=?", walletID).First(w).Error; err != nil {
			return err
		}

		w.Balance = w.Balance.Add(amount)
		if err := tx.Where("wallet_id=?", w.WalletID).Save(w).Error; err != nil {
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
		err := tx.Where("wallet_id IN(?)", []string{fromWalletID, toWalletID}).Find(&wallets).Error
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
			return fmt.Errorf("can't find wallet id %s", fromWalletID)
		}

		if to == nil {
			return fmt.Errorf("can't find wallet id %s", toWalletID)
		}

		if from.Balance.LessThan(amount) {
			return ErrBalanceNotEnough
		}

		from.Balance = from.Balance.Sub(amount)
		to.Balance = to.Balance.Add(amount)

		if err := tx.Where("wallet_id=?", from.WalletID).Save(from).Error; err != nil {
			return err
		}

		if err := tx.Where("wallet_id=?", to.WalletID).Save(to).Error; err != nil {
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
	return dao.db.Transaction(func(tx *gorm.DB) error {
		var count int64
		err := tx.Table(model.Wallet{}.TableName()).Where("wallet_id=?", walletID).Count(&count).Error
		if errors.Is(err, sql.ErrNoRows) || count == 0 {
			return ErrWalletNotExist
		}

		if err != nil {
			return err
		}

		return tx.Where("wallet_id=?", walletID).Delete(&model.Wallet{}).Error
	})
}
