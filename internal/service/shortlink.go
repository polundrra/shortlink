package service

import (
	"context"
	"github.com/polundrra/shortlink/internal/repo"
)

type linkService struct {
	repo repo.LinkRepo
}

func (l linkService) GetLongLink(ctx context.Context, code string) (string, error) {
	url, err := l.repo.GetLongLinkByCode(ctx, code)
	if err != nil {
		return "", err
	}

	if url == "" {
		return "", ErrLongLinkNotFound
	}

	return url, nil
}

func (l linkService) CreateShortLink(ctx context.Context, url, customEnd string) (string, error) {
	isCustom := false
	if customEnd != "" {
		isCustom = true

		exists, err := l.repo.IsCodeExists(ctx, customEnd)
		if err != nil {
			return "", err
		}

		if exists {
			exUrl, err := l.repo.GetLongLinkByCode(ctx, customEnd)
			if err != nil {
				return "", err
			}

			if url == exUrl {
				return customEnd, nil
			}

			return "", ErrCodeConflict
		}

		return customEnd, l.repo.SetLink(ctx, url, customEnd, isCustom)
	}

	code, err := l.repo.GetCodeByLongLink(ctx, url)
	if err != nil {
		return "", err
	}

	if code != "" {
		return code, nil
	}

	code, err = l.genCode(ctx)
	if err != nil {
		return "", err
	}

	if err := l.repo.SetLink(ctx, url, code, isCustom); err != nil {
		return "", err
	}

	return code, nil
}

func (l *linkService) genCode(ctx context.Context) (string, error) {
	exists := true
	var code string
	for exists {
		seq, err := l.repo.GetNextSeq(ctx)
		if err != nil {
			return "", err
		}

		code = toBase62(seq)
		exists, err = l.repo.IsCodeExists(ctx, code)
		if err != nil {
			return "", err
		}
	}

	return code, nil
}

func toBase62(n uint64) string {
	digits := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	length := uint64(len(digits))

	if n == 0 {
		return string(digits[0])
	}

	var res string
	for n > 0 {
		res = string(digits[n % length]) + res
		n = n / length
	}

	return res
}




