package service

import (
	"context"
	"github.com/polundrra/shortlink/internal/repo"
	"time"
)

type Service interface {
	GetLongLink(ctx context.Context, code string) (string, error)
	CreateShortLink(ctx context.Context, url, customEnd string) (string, error)
}

func New(opts repo.Opts, repo repo.LinkRepo) Service {
	return &linkService{
		repo: repo,
		timeout: time.Duration(opts.Timeout) * time.Second,
	}
}

