package os_saver

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"

	"project-evredika/internal/storage/data_saver"
)

type saver struct{}

func (s *saver) CreateData(_ context.Context, data *data_saver.Data) (err error) {
	var file *os.File
	if file, err = os.OpenFile(data.Bucket+data.Key, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0777); err != nil {
		return
	}

	if _, err = file.Write(data.B); err != nil {
		return
	}

	return nil
}

func (s *saver) UpdateData(_ context.Context, data *data_saver.Data) (err error) {
	var file *os.File
	if file, err = os.OpenFile(data.Bucket+data.Key, os.O_WRONLY, 0777); err != nil {
		return
	}

	if _, err = file.Write(data.B); err != nil {
		return
	}

	return nil
}

func (s *saver) ReadData(_ context.Context, info *data_saver.Metadata) (data []byte, err error) {
	return ioutil.ReadFile(info.Bucket + info.Key)
}

func (s *saver) DeleteData(_ context.Context, info *data_saver.Metadata) (err error) {
	return os.Remove(info.Bucket + info.Key)
}

func (s *saver) ListData(ctx context.Context, info *data_saver.Metadata) (data []*data_saver.Data, err error) {
	var fileInfo []os.FileInfo
	if fileInfo, err = ioutil.ReadDir(info.Bucket + info.Key); err != nil {
		return
	}

	for _, file := range fileInfo {
		if file.IsDir() {
			continue
		}

		key := filepath.Base(file.Name())

		var fileData []byte
		if fileData, err = s.ReadData(ctx, &data_saver.Metadata{
			Key:    key,
			Bucket: info.Bucket,
		}); err != nil {
			return
		}

		data = append(data, &data_saver.Data{
			Metadata: data_saver.Metadata{Key: key},
			B:        fileData,
		})
	}

	return data, nil
}

func (s *saver) Initiate(_ context.Context, bucket string) {
	_ = os.Mkdir(bucket, 0777)
}

func (s *saver) validateDataExist(info data_saver.Metadata) (err error) {
	var file *os.File
	file, err = os.Open(info.Bucket + info.Key)
	if err == nil {
		file.Close()
		return fmt.Errorf("file with key '%s' already exists", info.Bucket+info.Key)
	}

	if pErr, ok := err.(*os.PathError); ok && pErr.Err == syscall.ENOENT {
		return nil
	}

	return err
}

// NewOsSaver ...
func NewOsSaver() data_saver.DataSaver { return &saver{} }
