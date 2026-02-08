package flags

import (
	"errors"
)

type FileStoragePath string

func (f FileStoragePath) String() string {
	return string(f)
}

func (f *FileStoragePath) Set(s string) error {
	if s == "" {
		return errors.New("file storage path is empty")
	}
	*f = FileStoragePath(s)
	return nil
}
