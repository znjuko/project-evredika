package middlewares

import (
	"context"
	"time"

	"project-evredika/internal/storage/data_saver"
)

type storageLogger struct {
	logger logger

	next dataSaver
}

func (l *storageLogger) Initiate(ctx context.Context, bucket string) {
	l.next.Initiate(ctx, bucket)
}

func (l *storageLogger) CreateData(ctx context.Context, data *data_saver.Data) (err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "CreateData",
			"request", data,
			"response", err,
			"elapsed", time.Since(begin),
		)
	}(time.Now())
	return l.next.CreateData(ctx, data)
}

func (l *storageLogger) UpdateData(ctx context.Context, data *data_saver.Data) (err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "UpdateData",
			"request", data,
			"response", err,
			"elapsed", time.Since(begin),
		)
	}(time.Now())
	return l.next.UpdateData(ctx, data)
}

func (l *storageLogger) ReadData(ctx context.Context, info *data_saver.Metadata) (data []byte, err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "ReadData",
			"request", info,
			"response", err,
			"elapsed", time.Since(begin),
		)
	}(time.Now())
	return l.next.ReadData(ctx, info)
}

func (l *storageLogger) DeleteData(ctx context.Context, info *data_saver.Metadata) (err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "ReadData",
			"request", info,
			"response", err,
			"elapsed", time.Since(begin),
		)
	}(time.Now())
	return l.next.DeleteData(ctx, info)
}

func (l *storageLogger) ListData(ctx context.Context, info *data_saver.Metadata) (data []*data_saver.Data, err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "ReadData",
			"request", info,
			"response", err,
			"elapsed", time.Since(begin),
		)
	}(time.Now())
	return l.next.ListData(ctx, info)
}

// NewStorageLogger ...
func NewStorageLogger(
	logger logger,

	next dataSaver,
) data_saver.DataSaver {
	return &storageLogger{
		logger: logger,
		next:   next,
	}
}
