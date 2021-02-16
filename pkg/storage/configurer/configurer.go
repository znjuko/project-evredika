package configurer

import (
	"context"

	"project-evredika/internal/storage/data_saver"
)

const (
	OS = "OS"
	S3 = "s3"
)

type CfgHandler func(ctx context.Context, bucket string) (st data_saver.DataSaver, err error)

type StorageConfigurer interface {
	Configure(ctx context.Context, strType, bucket string) (st data_saver.DataSaver, err error)
}

type storageConfigurer struct {
	handlers map[string]CfgHandler
}

func (s *storageConfigurer) Configure(ctx context.Context, strType, bucket string) (st data_saver.DataSaver, err error) {
	handler := s.handlers[strType]
	return handler(ctx, bucket)
}

// NewStorageConfigurer ..
func NewStorageConfigurer(handlers map[string]CfgHandler) StorageConfigurer {
	return &storageConfigurer{handlers: handlers}
}
