package storages

import "io"

type Storage interface {
	Exists(name string) (bool, error)
	Read(name string) (io.ReadCloser, error)
	Write(name string) (io.WriteCloser, error)
}

type Initializer interface {
	Initialize() error
	Deinitialize()
}

type StorageInitializer interface {
	Storage
	Initializer
}
