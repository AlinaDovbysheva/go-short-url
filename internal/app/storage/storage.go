package storage

import (
	"encoding/json"
	"errors"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"sync"

	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
)

func NewInMap() DBurl {
	return &InMap{mapURL: map[string]string{}, mUid: []mapURLUid{}}
}

type InMap struct {
	mapURL map[string]string
	mUid   []mapURLUid
	mutex  sync.Mutex
}

type mapURLUid struct {
	Uid      string `json:"-"`
	URLShort string `json:"short_url"`
	URL      string `json:"original_url"`
}

func (m *InMap) Close() error {
	return nil
}

func (m *InMap) PingDB() error {
	return nil
}

func (m *InMap) GetAllURLUid(UID string) ([]byte, error) {
	var mUid []mapURLUid
	for _, v := range m.mUid {
		if v.Uid == UID {
			mUid = append(mUid, mapURLUid{v.Uid, v.URLShort, v.URL})
		}
	}
	if len(mUid) < 1 {
		return nil, errors.New("urls is absent in db")
	}
	data, _ := json.Marshal(mUid)
	return data, nil
}

func (m *InMap) GetURL(shortURL string, UID string) (string, error) {
	m.mutex.Lock()
	sID := m.mapURL[shortURL]
	m.mutex.Unlock()
	if sID == "" {
		return "", errors.New("id is absent in db")
	}
	return sID, nil
}

func (m *InMap) PutURL(inputURL string, UID string) (string, []byte, error) {
	id := ""
	for k, v := range m.mapURL {
		if v == inputURL {
			id = k
		}
	}
	if id == "" {
		id = util.RandStringBytes(7)
		m.mutex.Lock()
		m.mapURL[id] = inputURL
		m.mutex.Unlock()
	}

	// save UID and Url as history query
	m.mUid = append(m.mUid, mapURLUid{UID, app.BaseURL + `/` + id, inputURL})

	d := util.StrtoJSON(app.BaseURL + `/` + id)
	return id, d, nil
}

func (m *InMap) PutURLArray(inputURLJSON []byte, UID string) ([]byte, error) {
	return nil, nil
}
func (m *InMap) DelURLArray(inputURLJSON []byte, UID string) error {
	return nil
}
