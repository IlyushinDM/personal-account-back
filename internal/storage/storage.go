// Package storage ..
package storage

import (
	"context"
	"io"
)

// FileStorage определяет интерфейс для работы с файловым хранилищем.
type FileStorage interface {
	// Upload загружает файл в хранилище.
	// Возвращает уникальный ключ объекта и ошибку.
	Upload(ctx context.Context, file io.Reader, size int64, contentType string, objectKey string) error

	// Download скачивает файл из хранилища по его ключу.
	// Возвращает объект файла (io.ReadCloser) и ошибку.
	Download(ctx context.Context, objectKey string) (io.ReadCloser, error)
}
