package audit

import (
	"encoding/json"
	"os"
	"sync"
)

// FileObserver writes audit events into file.
type FileObserver struct {
	file *os.File
	mu   sync.Mutex
}

// NewFileObserver creates new file audit observer.
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

// Process writes audit event into file.
func (f *FileObserver) Process(event Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	f.mu.Lock()
	defer f.mu.Unlock()

	_, err = f.file.Write(data)
	return err
}

// Close closes audit log file.
func (f *FileObserver) Close() error {
	return f.file.Close()
}
