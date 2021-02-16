package data_saver

import (
	"context"
)

// DataSaver is common interface for data transitions between server-side and store
type DataSaver interface {
	CreateData(ctx context.Context, data *Data) (err error)
	UpdateData(ctx context.Context, data *Data) (err error)
	ReadData(ctx context.Context, info *Metadata) (data []byte, err error)
	DeleteData(ctx context.Context, info *Metadata) (err error)
	ListData(ctx context.Context, info *Metadata) (data []*Data, err error)
	Initiate(ctx context.Context, bucket string)
}

type Data struct {
	Metadata
	B []byte
}

type Metadata struct {
	Key, Bucket string
}
