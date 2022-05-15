package storage

type DBurl interface {
	GetURL(shortURL string) (string, error)
	PutURL(inputURL string) (string, error)
}
