package configurer

import (
	"context"

	"project-evredika/internal/storage/data_saver"
	"project-evredika/internal/storage/data_saver/os_saver"
)

func CreateOSStorage(ctx context.Context, bucket string) (st data_saver.DataSaver, err error) {
	st = os_saver.NewOsSaver()
	st.Initiate(ctx, bucket)
	return
}
