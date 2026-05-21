package flags

import (
	"errors"
)

// FileStoragePath represents path to metric storage file.
type FileStoragePath string

// String returns file storage path.
func (f FileStoragePath) String() string {
	return string(f)
}

// Set sets file storage path value.
func (f *FileStoragePath) Set(s string) error {
	if s == "" {
		return errors.New("file storage path is empty")
	}
	*f = FileStoragePath(s)
	return nil
}
