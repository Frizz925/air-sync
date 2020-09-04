package services

type Initializer interface {
	Initialize() error
	Deinitialize()
}
