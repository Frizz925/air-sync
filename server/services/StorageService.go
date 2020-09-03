package services

import (
	"air-sync/storages"
	"context"
)

type StorageService struct {
	storage storages.Storage
}

var _ Service = (*StorageService)(nil)

func NewStorageService(ctx context.Context, bucketName string) *StorageService {
	return &StorageService{
		storage: storages.NewGoogleCloudStorage(ctx, bucketName),
	}
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
