package services

type Service interface {
	Initialize() error
	Deinitialize()
}
