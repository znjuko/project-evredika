package transaction

import (
	"context"

	"project-evredika/internal/storage/data_saver"
)

type dataSaver interface {
	CreateData(ctx context.Context, data *data_saver.Data) (err error)
	UpdateData(ctx context.Context, data *data_saver.Data) (err error)
	ReadData(ctx context.Context, info *data_saver.Metadata) (data []byte, err error)
	DeleteData(ctx context.Context, info *data_saver.Metadata) (err error)
	ListData(ctx context.Context, info *data_saver.Metadata) (data []*data_saver.Data, err error)
}
