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

func (p *postgres) SetLink(ctx context.Context, url, code string, isCustom bool) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	if _, err := p.pool.ExecEx(ctx, "insert into link(url, code, is_custom) values ($1, $2, $3)", nil, url, code, isCustom); err != nil {
		return fmt.Errorf("error insert url %s, code %s: %w", url, code, err)
	}

	return nil
}

func (p *postgres) GetLongLinkByCode(ctx context.Context, code string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var url string
	if err := p.pool.QueryRowEx(ctx, "select url from link where code=$1", nil, code).Scan(&url); err != nil {
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
	if err := p.pool.QueryRowEx(ctx, "select code from link where url=$1 and is_custom=$2", nil, url, false).Scan(&code); err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}

		return "", fmt.Errorf("error get code by url %s: %w", url, err)
	}

	return code, nil
}

func (p *postgres) GetNextSeq(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var nextVal uint64
	if err := p.pool.QueryRowEx(ctx, "select nextval('seq')", nil).Scan(&nextVal); err != nil {
		return 0, fmt.Errorf("error get next val %w", err)
	}

	return nextVal, nil
}

func (p *postgres) IsCodeExists(ctx context.Context, code string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	var c int
	if err := p.pool.QueryRowEx(ctx, "select count(*) from link where code=$1", nil, code).Scan(&c); err != nil {
		return false, fmt.Errorf("error check if code exists %w", err)
	}

	return c != 0, nil
}



