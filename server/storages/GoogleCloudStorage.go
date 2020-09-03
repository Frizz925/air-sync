package storages

import (
	"context"
	"errors"
	"io"

	log "github.com/sirupsen/logrus"

	. "cloud.google.com/go/storage"
)

type GoogleCloudStorage struct {
	context    context.Context
	bucket     *BucketHandle
	bucketName string
}

var _ StorageInitializer = (*GoogleCloudStorage)(nil)

func NewGoogleCloudStorage(ctx context.Context, bucketName string) *GoogleCloudStorage {
	return &GoogleCloudStorage{
		context:    ctx,
		bucketName: bucketName,
	}
}

func (s *GoogleCloudStorage) Initialize() error {
	client, err := NewClient(s.context)
	if err != nil {
		return err
	}
	s.bucket = client.Bucket(s.bucketName)
	log.Infof("Using Google Cloud Storage: %s", s.bucketName)
	return nil
}

func (s *GoogleCloudStorage) Deinitialize() {
}

func (s *GoogleCloudStorage) Exists(name string) (bool, error) {
	_, err := s.bucket.Object(name).Attrs(s.context)
	if errors.Is(err, ErrObjectNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func (s *GoogleCloudStorage) Read(name string) (io.ReadCloser, error) {
	return s.bucket.Object(name).NewReader(s.context)
}

func (s *GoogleCloudStorage) Write(name string) (io.WriteCloser, error) {
	return s.bucket.Object(name).NewWriter(s.context), nil
}
