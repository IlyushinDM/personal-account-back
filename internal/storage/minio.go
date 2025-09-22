package storage

import (
	"context"
	"fmt"
	"io"

	"lk/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinIOClient реализует интерфейс FileStorage для MinIO.
type MinIOClient struct {
	client     *minio.Client
	bucketName string
}

// NewMinIOClient создает и настраивает нового клиента MinIO.
func NewMinIOClient(cfg config.MinioConfig) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("could not create minio client: %w", err)
	}

	// Проверяем, существует ли бакет, и создаем его, если нет.
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, cfg.BucketName)
	if err != nil {
		return nil, fmt.Errorf("could not check if bucket exists: %w", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("could not create bucket: %w", err)
		}
	}

	return &MinIOClient{
		client:     client,
		bucketName: cfg.BucketName,
	}, nil
}

// Upload загружает файл в MinIO.
func (m *MinIOClient) Upload(ctx context.Context, file io.Reader, size int64, contentType string, objectKey string) error {
	_, err := m.client.PutObject(ctx, m.bucketName, objectKey, file, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// Download скачивает файл из MinIO.
func (m *MinIOClient) Download(ctx context.Context, objectKey string) (io.ReadCloser, error) {
	object, err := m.client.GetObject(ctx, m.bucketName, objectKey, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return object, nil
}
