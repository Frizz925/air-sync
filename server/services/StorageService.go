package services

import (
	"air-sync/storages"
	"context"
)

type StorageOptions struct {
	BucketName string
	UploadsDir string
}

type StorageService struct {
	fileStorage  *storages.FileStorage
	cloudStorage *storages.GoogleCloudStorage
}

var _ Service = (*StorageService)(nil)

func NewStorageService(ctx context.Context, opts StorageOptions) *StorageService {
	return &StorageService{
		fileStorage:  storages.NewFileStorage(opts.UploadsDir),
		cloudStorage: storages.NewGoogleCloudStorage(ctx, opts.BucketName),
	}
}

func (s *StorageService) Initialize() error {
	if err := s.fileStorage.Initialize(); err != nil {
		return err
	}
	if err := s.cloudStorage.Initialize(); err != nil {
		return err
	}
	return nil
}

func (s *StorageService) Deinitialize() {
	s.fileStorage.Deinitialize()
	s.cloudStorage.Deinitialize()
}

func (s *StorageService) FileStorage() *storages.FileStorage {
	return s.fileStorage
}

func (s *StorageService) CloudStorage() *storages.GoogleCloudStorage {
	return s.cloudStorage
}
