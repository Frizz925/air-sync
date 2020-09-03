package storages

import (
	"io"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type FileStorage struct {
	dir    string
	absDir string
}

var _ StorageInitializer = (*FileStorage)(nil)

func NewFileStorage(dir string) *FileStorage {
	return &FileStorage{
		dir: dir,
	}
}

func (s *FileStorage) Initialize() error {
	path, err := filepath.Abs(s.dir)
	if err != nil {
		return err
	}
	log.Infof("Using file storage: %s", path)
	// Create directory if not exists
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		log.Infof("Creating directory: %s", path)
		if err := os.Mkdir(path, 0755); err != nil {
			return err
		}
	}
	s.absDir = path
	return nil
}

func (s *FileStorage) Deinitialize() {
	// Do nothing
}

func (s *FileStorage) Exists(name string) (bool, error) {
	_, err := os.Stat(s.getPath(name))
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *FileStorage) Read(name string) (io.ReadCloser, error) {
	return os.Open(s.getPath(name))
}

func (s *FileStorage) Write(name string) (io.WriteCloser, error) {
	return os.Create(s.getPath(name))
}

func (s *FileStorage) getPath(name string) string {
	return filepath.Join(s.absDir, filepath.Clean(name))
}
