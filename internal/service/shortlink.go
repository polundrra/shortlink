package service

import (
	"context"
	"github.com/jackc/pgx"
	"github.com/polundrra/shortlink/internal/repo"
	"regexp"
	"strings"
	"time"
)

type linkService struct {
	repo repo.LinkRepo
	timeout time.Duration
}

func (l linkService) GetLongLink(ctx context.Context, code string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	url, err := l.repo.GetLongLinkByCode(ctx, code)
	if err != nil {
		return "", err
	}
	if url == "" {
		return "", ErrLongLinkNotFound
	}
	return url, nil
}

func (l linkService) GetShortLink(ctx context.Context, url, customEnd string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	code, err := l.repo.GetCodeByLongLink(ctx, url)
	if err != nil {
		return "", err
	}

	if code != "" {
		if customEnd == "" {
			return code, nil
		}
		return code, ErrShortLinkExists
	}
	if code == "no rows" {
		if err := l.repo.AddLongLink(ctx, url); err != nil {
			return "", err
		}
	}

	if customEnd != ""{
		if !validateEnding(customEnd) {
			return "", ErrInvalidEnding
		}
		if err := l.repo.SetShortLink(ctx, url, customEnd); err != nil {
			notUnique:= pgx.PgError{Code: "23505"}
			if err == notUnique {
				return "", ErrShortLinkExists
			}
			return "", err
		}
		return customEnd, nil
	}

	nextVal, err := l.repo.GetNextSeq(ctx)
	if err != nil {
		return "", err
	}

	if err := l.repo.SetShortLink(ctx, url, toBase62(nextVal)); err != nil {
		notUnique:= pgx.PgError{Code: "23505"}
		if err == notUnique {
			for err == notUnique {
				nextVal, err = l.repo.GetNextSeq(ctx)
				if err != nil {
					return "", err
				}
				err = l.repo.SetShortLink(ctx, url, toBase62(nextVal))
			}
			if err != nil {
				return "", err
			}
		}
		return "", err
	}
	return toBase62(nextVal), nil
}

func toBase62(n uint64) string {
	digits := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := uint64(len(digits))
	if n == 0 {
		return string(digits[0])
	}
	var sb strings.Builder
	for ; n > 0; n = n / length {
		sb.WriteByte(digits[n % length])
	}
	return sb.String()
}

func validateEnding(ending string) bool {
	regex:= regexp.MustCompile("^[a-zA-Z0-9-_]+&")
	return regex.MatchString(ending)
}




