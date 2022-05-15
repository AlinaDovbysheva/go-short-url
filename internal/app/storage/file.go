package storage

import (
	"encoding/json"
	"os"
)

func NewInFile() DBurl {
	fileName := "events.log"
	defer os.Remove(fileName)

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil
	}
	return &InFile{file: file, encoder: json.NewEncoder(file)}

}

type Event struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type InFile struct {
	file    *os.File
	encoder *json.Encoder
}

func (p *InFile) GetURL(shortURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *InFile) PutURL(inputURL string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewInFile2(fileName string) (*InFile, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}
	return &InFile{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}
func (p *InFile) WriteEvent(event *Event) error {
	return p.encoder.Encode(&event)
}
func (p *InFile) Close() error {
	return p.file.Close()
}
