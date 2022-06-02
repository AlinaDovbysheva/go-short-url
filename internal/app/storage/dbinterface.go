package storage

type DBurl interface {
	GetURL(shortURL string) (string, error)
	PutURL(inputURL string, UID string) (string, []byte, error)
	GetAllURLUid(UID string) ([]byte, error)

	// in- array(json) of original url out-array(json) of short url
	PutURLArray(inputURLJSON []byte, UID string) ([]byte, error)
	Close() error  //нужна только для db
	PingDB() error //нужна только для db
}
