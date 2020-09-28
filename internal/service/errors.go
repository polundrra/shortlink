package service

type LinkError string

func (e LinkError) Error() string {
	return string(e)
}

const (
	ErrLongLinkNotFound LinkError = "Long link not found, check correctness of short link"
	ErrCodeConflict LinkError = "Short link already exists"
	ErrInvalidCode LinkError = "Custom link must contain only alphanumeric characters, hyphens, underscores with length 1-32"
)