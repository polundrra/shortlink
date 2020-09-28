package repo

import (
	"context"
	"github.com/jackc/pgx"
	"time"
)

type LinkRepo interface {
	SetLink(ctx context.Context, url, code string, isCustom bool) error
	GetLongLinkByCode(ctx context.Context, code string) (string, error)
	GetCodeByLongLink(ctx context.Context, url string) (string, error)
	GetNextSeq(ctx context.Context) (uint64, error)
	IsCodeExists(ctx context.Context, code string) (bool, error)
}

type Opts struct {
	Host string
	Port uint16
	Database string
	User string
	Password string
	Timeout int
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
		timeout: time.Duration(opts.Timeout) * time.Second,
	}
	return &repo, nil
}