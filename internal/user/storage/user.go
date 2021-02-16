package storage

import (
	"context"
	"encoding/json"
	"fmt"

	"project-evredika/internal/storage/data_saver"
	"project-evredika/internal/user"
	v1 "project-evredika/pkg/api/v1"
)

type userStorage struct {
	storage    storage
	indexer    dataIndexer
	indexMaker indexMaker
	logger     logger
	bucket     string
}

func (s *userStorage) CreateUser(ctx context.Context, user *v1.User) (err error) {
	d := &data_saver.Data{
		Metadata: data_saver.Metadata{
			Key:    s.indexMaker.CreateIndex(user.ID),
			Bucket: s.bucket,
		},
	}

	s.storage.StartTransaction(d.Key)
	defer s.storage.StopTransaction(d.Key)

	s.logger.Debug("created data key", d.Key, "bucket", d.Bucket)

	if d.B, err = json.Marshal(user); err != nil {
		return
	}

	if err = s.storage.CreateData(ctx, d); err != nil {
		return
	}

	return nil
}

func (s *userStorage) DeleteUser(ctx context.Context, ID string) (err error) {
	i := &data_saver.Metadata{
		Key:    s.indexMaker.CreateIndex(ID),
		Bucket: s.bucket,
	}

	s.storage.StartTransaction(i.Key)
	defer s.storage.StopTransaction(i.Key)

	s.logger.Debug("created data key", i.Key, "bucket", i.Bucket)

	if err = s.storage.DeleteData(ctx, i); err != nil {
		return
	}

	return nil
}

func (s *userStorage) UpdateUser(ctx context.Context, user *v1.User) (err error) {
	d := &data_saver.Data{
		Metadata: data_saver.Metadata{
			Key:    s.indexMaker.CreateIndex(user.ID),
			Bucket: s.bucket,
		},
	}

	s.storage.StartTransaction(d.Key)
	defer s.storage.StopTransaction(d.Key)

	s.logger.Debug("created data key", d.Key, "bucket", d.Bucket)

	if d.B, err = json.Marshal(user); err != nil {
		return
	}

	if err = s.storage.UpdateData(ctx, d); err != nil {
		return
	}

	return nil
}

func (s *userStorage) GetUser(_ context.Context, ID string) (user *v1.User, err error) {
	var exist bool
	if user, exist = s.indexer.GetUser(ID); !exist {
		return nil, fmt.Errorf("user %v not found", ID)
	}

	return user, nil
}

func (s *userStorage) ListUsers(_ context.Context, skip, limit int) (users []*v1.User) {
	return s.indexer.ListUsers(skip, limit)
}

// NewUserStorage ...
func NewUserStorage(
	storage storage,
	indexer dataIndexer,
	indexMaker indexMaker,
	logger logger,
	bucket string,
) user.Storage {
	return &userStorage{
		storage:    storage,
		indexer:    indexer,
		indexMaker: indexMaker,
		logger:     logger,
		bucket:     bucket,
	}
}
