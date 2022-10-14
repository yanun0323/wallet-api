package service

import (
	"context"
	"wallet-api/internal/domain"
	"wallet-api/internal/repository"

	"github.com/labstack/echo/v4"
)

type Service struct {
	l    echo.Logger
	ctx  context.Context
	repo domain.Repository
}

func New(ctx context.Context, l echo.Logger) (*Service, error) {
	repo, err := repository.New()
	if err != nil {
		return nil, err
	}
	return &Service{
		ctx:  ctx,
		repo: repo,
		l:    l,
	}, nil
}
