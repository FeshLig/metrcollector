package audit

import (
	"encoding/json"
	"os"
)

type FileObserver struct {
	file *os.File
}

func NewFileObserver(path string) (*FileObserver, error) {
	file, err := os.OpenFile(
		path,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		return nil, err
	}

	return &FileObserver{
		file: file,
	}, nil
}

func (f *FileObserver) Process(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = f.file.Write(data)
	return err
}

func (f *FileObserver) Close() error {
	return f.file.Close()
}
