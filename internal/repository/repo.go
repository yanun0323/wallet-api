package repository

import (
	"wallet-api/internal/domain"
	"wallet-api/internal/repository/sqlite"

	"github.com/spf13/viper"
)

type Repo struct {
	*sqlite.SqliteDao
}

func New() (domain.Repository, error) {
	sqliteDao, err := sqlite.NewSqliteDao(viper.GetString("sqlite.db"))
	if err != nil {
		return nil, err
	}
	return &Repo{
		SqliteDao: sqliteDao,
	}, nil
}
