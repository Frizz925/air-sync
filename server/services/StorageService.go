package services

import (
	"air-sync/storages"
	"context"
)

type StorageOptions struct {
	StorageMode string
	BucketName  string
	UploadsDir  string
}

type StorageService struct {
	storage storages.Storage
}

var _ Service = (*StorageService)(nil)

func NewStorageService(ctx context.Context, opts StorageOptions) *StorageService {
	fileStorage := storages.NewFileStorage(opts.UploadsDir)
	cloudStorage := storages.NewGoogleCloudStorage(ctx, opts.BucketName)
	service := &StorageService{}

	switch opts.StorageMode {
	case "Cache":
		service.storage = storages.NewCacheStorage(
			fileStorage,
			cloudStorage,
		)
	case "GoogleCloudStorage":
		service.storage = cloudStorage
	default:
		service.storage = fileStorage
	}

	return service
}

func (s *StorageService) Initialize() error {
	if v, ok := s.storage.(storages.Initializer); ok {
		return v.Initialize()
	}
	return nil
}

func (s *StorageService) Deinitialize() {
	if v, ok := s.storage.(storages.Initializer); ok {
		v.Deinitialize()
	}
}

func (s *StorageService) Storage() storages.Storage {
	return s.storage
}
