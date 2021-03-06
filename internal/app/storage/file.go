package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"io"
	"log"
	"os"
	"sync"

	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
)

type Event struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type InFile struct {
	fileName string
	mapURL   map[string]string
	mutex    *sync.Mutex
}

func NewInFile(fileName string) DBurl {
	rf, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer rf.Close()

	m := map[string]string{}

	dec := json.NewDecoder(rf)
	for {
		emp := Event{}
		if err := dec.Decode(&emp); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		m[emp.ID] = emp.URL
	}
	fmt.Println("read from file " + fileName + " to map:")
	fmt.Println(m)
	return &InFile{fileName: fileName, mapURL: m}
}

func (m *InFile) Close() error {
	return nil
}

func (m *InFile) PingDB() error {
	return nil
}

func (m *InFile) GetAllURLUid(UID string) ([]byte, error) {
	return nil, nil
}
func (m *InFile) PutURLArray(inputURLJSON []byte, UID string) ([]byte, error) {
	return nil, nil
}
func (m *InFile) DelURLArray(inputURLJSON []byte, UID string) error {
	return nil
}

func (m *InFile) GetURL(shortURL string, UID string) (string, error) {
	sID := m.mapURL[shortURL]
	if sID == "" {
		return "", errors.New("id is absent in db")
	}
	return sID, nil
}

func (m *InFile) PutURL(inputURL string, UID string) (string, []byte, error) {
	id := ""
	for k, v := range m.mapURL {
		if v == inputURL {
			id = k
		}
	}

	if id == "" {
		id = util.RandStringBytes(7)

		//p.mutex.Lock()
		m.mapURL[id] = inputURL
		//p.mutex.Unlock()

		//------write to file
		event := Event{id, inputURL}
		wf, err := NewWFile(m.fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer wf.CloseWFile()
		//p.mutex.Lock()
		if err := wf.WriteEvent(event); err != nil {
			log.Fatal(err)
		}
		//p.mutex.Unlock()
	}
	d := util.StrtoJSON(app.BaseURL + `/` + id)
	return id, d, nil
}

type WFile struct {
	file    *os.File
	encoder *json.Encoder
}

func NewWFile(fileName string) (*WFile, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &WFile{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (m *WFile) WriteEvent(event Event) error {
	return m.encoder.Encode(&event)
}

func (m *WFile) CloseWFile() error {
	return m.file.Close()
}

type RFile struct {
	file    *os.File
	decoder *json.Decoder
}

func NewRFile(fileName string) (*RFile, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return nil, err
	}
	return &RFile{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}
func (c *RFile) ReadEvent() (*Event, error) {
	event := &Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}
	return event, nil
}
func (c *RFile) CloseRFile() error {
	return c.file.Close()
}
