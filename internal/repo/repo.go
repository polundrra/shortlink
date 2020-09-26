package repo

import (
	"context"
	"github.com/jackc/pgx"
	"time"
)

type LinkRepo interface {
	AddLongLink(url string, ctx context.Context) error
	SetShortLink(url, code string, ctx context.Context) error
	GetLongLinkByCode(code string, ctx context.Context) (string, error)
	GetNextSeq(ctx context.Context) (uint64, error)
}

type Opts struct {
	Host string
	Port uint16
	Database string
	User string
	Password string
	Timeout time.Duration
}

func New(opts Opts) (LinkRepo, error) {
	ConnConfig := pgx.ConnConfig{
		Host: opts.Host,
		Port: opts.Port,
		Database: opts.Database,
		User: opts.User,
		Password: opts.Password,
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: ConnConfig,
	})
	if err != nil {
		return nil, err
	}
	repo := postgres{
		pool: pool,
		timeout: opts.Timeout,
	}
	return &repo, nil
}