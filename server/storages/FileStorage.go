package storages

import (
	"io"
	"os"
	"path/filepath"
)

type FileStorage struct {
	dir string
}

var _ Storage = (*FileStorage)(nil)

func NewFileStorage(dir string) *FileStorage {
	return &FileStorage{dir}
}

func (s *FileStorage) Exists(name string) (bool, error) {
	path, err := s.getPath(name)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *FileStorage) Read(name string) (io.ReadCloser, error) {
	return s.getFile(name)
}

func (s *FileStorage) Write(name string) (io.WriteCloser, error) {
	return s.getFile(name)
}

func (s *FileStorage) getFile(name string) (*os.File, error) {
	path, err := s.getPath(name)
	if err != nil {
		return nil, err
	}
	return os.Open(path)
}

func (s *FileStorage) getPath(name string) (string, error) {
	path, err := filepath.Abs(name)
	if err != nil {
		return "", err
	}
	return filepath.Join(s.dir, path), nil
}
