package storages

import (
	"errors"
	"io"
)

var ErrObjectNotFound = errors.New("Object not found")

type CacheStorage struct {
	storages []Storage
}

type CacheWriteCloser struct {
	io.Writer
	closers []io.Closer
}

type CacheReadWriteCloser struct {
	*CacheWriteCloser
	io.ReadCloser
}

var (
	_ StorageInitializer = (*CacheStorage)(nil)
	_ io.WriteCloser     = (*CacheWriteCloser)(nil)
)

func NewCacheStorage(storages ...Storage) *CacheStorage {
	return &CacheStorage{storages}
}

func (s *CacheStorage) Initialize() error {
	for _, storage := range s.storages {
		init, ok := storage.(Initializer)
		if !ok {
			continue
		}
		if err := init.Initialize(); err != nil {
			return err
		}
	}
	return nil
}

func (s *CacheStorage) Deinitialize() {
	for _, storage := range s.storages {
		if init, ok := storage.(Initializer); ok {
			init.Deinitialize()
		}
	}
}

func (s *CacheStorage) Exists(name string) (bool, error) {
	for _, storage := range s.storages {
		exists, err := storage.Exists(name)
		if err != nil {
			return false, err
		} else if exists {
			return true, nil
		}
	}
	return false, nil
}

func (s *CacheStorage) Read(name string) (io.ReadCloser, error) {
	rwc := &CacheReadWriteCloser{}
	writers := make([]io.WriteCloser, 0)
	for _, storage := range s.storages {
		exists, err := storage.Exists(name)
		if err != nil {
			return nil, err
		} else if !exists {
			// If not exists, then add it as one of the writers to write during reads
			w, err := storage.Write(name)
			if err != nil {
				return nil, err
			}
			writers = append(writers, w)
		} else if rwc.ReadCloser == nil {
			// If exists but no reader is assigned, then use it as the main source for reads
			r, err := storage.Read(name)
			if err != nil {
				return nil, err
			}
			rwc.ReadCloser = r
		}
	}
	if rwc.ReadCloser == nil {
		return nil, ErrObjectNotFound
	}
	rwc.CacheWriteCloser = NewCacheWriteCloser(writers...)
	return rwc, nil
}

func (s *CacheStorage) Write(name string) (io.WriteCloser, error) {
	writers := make([]io.WriteCloser, len(s.storages))
	for idx, storage := range s.storages {
		w, err := storage.Write(name)
		if err != nil {
			return nil, err
		}
		writers[idx] = w
	}
	return NewCacheWriteCloser(writers...), nil
}

func (s *CacheStorage) Delete(name string) error {
	for _, storage := range s.storages {
		exists, err := storage.Exists(name)
		if err != nil {
			return err
		} else if !exists {
			continue
		}
		if err := storage.Delete(name); err != nil {
			return err
		}
	}
	return nil
}

func (rwc *CacheReadWriteCloser) Read(b []byte) (int, error) {
	n, err := rwc.ReadCloser.Read(b)
	if err != nil {
		return n, err
	}
	if _, err := rwc.Write(b[:n]); err != nil {
		return n, err
	}
	return n, nil
}

func (rwc *CacheReadWriteCloser) Close() error {
	if err := rwc.ReadCloser.Close(); err != nil {
		return err
	}
	if err := rwc.CacheWriteCloser.Close(); err != nil {
		return err
	}
	return nil
}

func NewCacheWriteCloser(writeClosers ...io.WriteCloser) *CacheWriteCloser {
	writers := make([]io.Writer, len(writeClosers))
	closers := make([]io.Closer, len(writeClosers))
	for idx, w := range writeClosers {
		writers[idx] = w
		closers[idx] = w
	}
	return &CacheWriteCloser{
		Writer:  io.MultiWriter(writers...),
		closers: closers,
	}
}

func (wc *CacheWriteCloser) Close() error {
	for _, closer := range wc.closers {
		if err := closer.Close(); err != nil {
			return err
		}
	}
	return nil
}
