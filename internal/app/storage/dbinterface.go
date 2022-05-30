package storage

type DBurl interface {
	GetURL(shortURL string) (string, error)
	PutURL(inputURL string, UID string) (string, error)
	GetAllURLUid(UID string) ([]byte, error)
	Close() error  //нужна только для db
	PingDB() error //нужна только для db
}
