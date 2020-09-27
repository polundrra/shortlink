package service

type LinkError string

func (e LinkError) Error() string {
	return string(e)
}

const (
	ErrLongLinkNotFound LinkError = "Long link not found, check correctness of short link"
	ErrShortLinkExists LinkError = "Short link already exists"
	ErrInvalidEnding LinkError = "Url ending must only contain alphanumeric characters, hyphens, underscores"
)