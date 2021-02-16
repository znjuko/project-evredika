package middlewares

import (
	"context"

	"project-evredika/internal/eventer/models"
	"project-evredika/internal/storage/data_saver"
)

type storageEventer struct {
	eventSender eventSender
	indexer     indexer

	next dataSaver
}

func (s *storageEventer) Initiate(ctx context.Context, bucket string) {
	s.next.Initiate(ctx, bucket)
}

func (s *storageEventer) CreateData(ctx context.Context, data *data_saver.Data) (err error) {
	if err = s.next.CreateData(ctx, data); err != nil {
		return
	}

	s.eventSender.SendEvent(&models.Event{
		Type: models.EventPut,
		Metadata: &models.Metadata{
			Data: data.B,
			Key:  s.indexer.RevertIndex(data.Key),
		},
	})

	return
}

func (s *storageEventer) UpdateData(ctx context.Context, data *data_saver.Data) (err error) {
	if err = s.next.UpdateData(ctx, data); err != nil {
		return
	}

	s.eventSender.SendEvent(&models.Event{
		Type: models.EventPut,
		Metadata: &models.Metadata{
			Data: data.B,
			Key:  s.indexer.RevertIndex(data.Key),
		},
	})

	return
}

func (s *storageEventer) ReadData(ctx context.Context, info *data_saver.Metadata) (data []byte, err error) {
	return s.next.ReadData(ctx, info)
}

func (s *storageEventer) DeleteData(ctx context.Context, info *data_saver.Metadata) (err error) {
	if err = s.next.DeleteData(ctx, info); err != nil {
		return
	}

	s.eventSender.SendEvent(&models.Event{
		Type: models.EventDelete,
		Metadata: &models.Metadata{
			Key: s.indexer.RevertIndex(info.Key),
		},
	})

	return
}

func (s *storageEventer) ListData(ctx context.Context, info *data_saver.Metadata) (data []*data_saver.Data, err error) {
	if data, err = s.next.ListData(ctx, info); err != nil {
		return
	}

	events := make([]*models.Event, len(data))
	for i, d := range data {
		events[i] = &models.Event{
			Type: models.EventPut,
			Metadata: &models.Metadata{
				Data: d.B,
				Key:  s.indexer.RevertIndex(d.Key),
			},
		}
	}

	s.eventSender.SendEvent(events...)

	return
}

// NewStorageEventer ...
func NewStorageEventer(
	eventSender eventSender,
	indexer indexer,
	next dataSaver,
) data_saver.DataSaver {
	return &storageEventer{
		eventSender: eventSender,
		indexer:     indexer,
		next:        next,
	}
}
