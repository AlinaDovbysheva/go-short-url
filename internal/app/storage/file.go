package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	"io"
	"log"
	"os"
)

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

type Event struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type InFile struct {
	fileName string
	mapURL   map[string]string
}

func (p *InFile) GetURL(shortURL string) (string, error) {
	sID := p.mapURL[shortURL]
	if sID == "" {
		return "", errors.New("id is absent in db")
	}
	return sID, nil
}

func (p *InFile) PutURL(inputURL string) (string, error) {
	id := ""
	for k, v := range p.mapURL {
		if v == inputURL {
			id = k
			fmt.Println(" __find url:")
			fmt.Println(id)
		}
	}

	if id == "" {
		id = util.RandStringBytes(7)
		p.mapURL[id] = inputURL

		//------write to file
		event := Event{id, inputURL}
		wf, err := NewWFile(p.fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer wf.Close()
		if err := wf.WriteEvent(event); err != nil {
			log.Fatal(err)
		}
	}

	return id, nil
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
func (p *WFile) WriteEvent(event Event) error {
	return p.encoder.Encode(&event)
}
func (p *WFile) Close() error {
	return p.file.Close()
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
func (c *RFile) Close() error {
	return c.file.Close()
}
