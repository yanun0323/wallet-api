package service

import (
	"context"
	"wallet-api/internal/domain"
	"wallet-api/internal/repository"

	"github.com/yanun0323/pkg/logs"
)

type Service struct {
	l    *logs.Logger
	ctx  context.Context
	repo domain.Repository
}

func New(ctx context.Context, l *logs.Logger) (*Service, error) {
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
