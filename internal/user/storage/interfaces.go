package storage

import (
	"context"

	"project-evredika/internal/storage/data_saver"
	v1 "project-evredika/pkg/api/v1"
)

type dataIndexer interface {
	GetUser(ID string) (user *v1.User, exist bool)
	ListUsers(skip, limit int) (users []*v1.User)
}

type indexMaker interface {
	CreateIndex(id string) (index string)
}

type storage interface {
	CreateData(ctx context.Context, data *data_saver.Data) (err error)
	UpdateData(ctx context.Context, data *data_saver.Data) (err error)
	DeleteData(ctx context.Context, info *data_saver.Metadata) (err error)
	StartTransaction(key string)
	StopTransaction(key string)
}

type logger interface {
	Error(args ...interface{})
	Debug(args ...interface{})
}