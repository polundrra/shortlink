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



func (p *postgres) AddLongLink(ctx context.Context, url string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.Exec("insert into link(url) values ($1)", url); err != nil {
		return fmt.Errorf("error insert url %s: %w", url, err)
	}
	return nil
}

func (p *postgres) SetShortLink(ctx context.Context, url, code string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.Exec("insert into link(code) values ($1) where url = '$2'", code, url); err != nil {
		return fmt.Errorf("error insert code %s: %w", code, err)
	}
	return nil
}

func (p *postgres) GetLongLinkByCode(ctx context.Context, code string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var url string
	if err := p.pool.QueryRow("select url from link where code = &1", code).Scan(&url); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("error get url by code %s: %w", code, err)
	}
	return url, nil
}

func (p *postgres) GetCodeByLongLink(ctx context.Context, url string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var code string
	if err := p.pool.QueryRow("select code from link where url = &1", url).Scan(&code); err != nil {
		if err == pgx.ErrNoRows {
			return "no rows", nil
		}
		return "", fmt.Errorf("error get code by url %s: %w", url, err)
	}
	return code, nil
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



