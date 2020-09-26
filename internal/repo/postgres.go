package repo

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"time"
)

type postgres struct {
	pool *pgx.ConnPool
	timeout time.Duration
}

func (p *postgres) AddLongLink(url string, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.Exec("insert into link(url) values ($1)", url); err != nil {
		return fmt.Errorf("error insert url %s: %w", url, err)
	}
	return nil
}

func (p *postgres) SetShortLink(url, code string, ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.Exec("insert into link(code) values ($1) where url = '$2'", code, url); err != nil {
		return fmt.Errorf("error insert code %s: %w", code, err)
	}
	return nil
}

func (p *postgres) GetLongLinkByCode(code string, ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var url string
	if err := p.pool.QueryRow("select url from link where code = &1", code).Scan(&url); err != nil {
		return "", fmt.Errorf("error get url by code %s: %w", code, err)
	}
	return url, nil
}

func (p *postgres) GetNextSeq(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var nextVal uint64
	if err := p.pool.QueryRow("select nextval('seq')").Scan(&nextVal); err != nil {
		return 0, fmt.Errorf("error get next val %w", err)
	}
	return nextVal, nil
}



