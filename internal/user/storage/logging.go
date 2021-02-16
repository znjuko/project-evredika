package storage

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"project-evredika/internal/user"
	v1 "project-evredika/pkg/api/v1"
)

type logging struct {
	logger *logrus.Logger

	next user.Storage
}

func (l *logging) CreateUser(ctx context.Context, user *v1.User) (err error) {
	defer func(begin time.Time) {
		l.logger.WithFields(logrus.Fields{
			"request": user,
			"error":   err,
		}).Info("CreateUser")
	}(time.Now())
	return l.next.CreateUser(ctx, user)
}

func (l *logging) DeleteUser(ctx context.Context, ID string) (err error) {
	defer func(begin time.Time) {
		l.logger.WithFields(logrus.Fields{
			"request": ID,
			"error":   err,
		}).Info("GetUser")
	}(time.Now())
	return l.next.DeleteUser(ctx, ID)
}

func (l *logging) UpdateUser(ctx context.Context, user *v1.User) (err error) {
	defer func(begin time.Time) {
		l.logger.WithFields(logrus.Fields{
			"request": user,
			"error":   err,
		}).Info("UpdateUser")
	}(time.Now())
	return l.next.UpdateUser(ctx, user)
}

func (l *logging) GetUser(ctx context.Context, ID string) (user *v1.User, err error) {
	defer func(begin time.Time) {
		l.logger.WithFields(logrus.Fields{
			"request":  ID,
			"response": user,
			"error":    err,
		}).Info("GetUser")
	}(time.Now())
	return l.next.GetUser(ctx, ID)
}

func (l *logging) ListUsers(ctx context.Context, skip, limit int) (users []*v1.User) {
	defer func(begin time.Time) {
		l.logger.WithFields(logrus.Fields{
			"skip":     skip,
			"limit":    limit,
			"response": users,
		}).Info("ListUsers")
	}(time.Now())
	return l.next.ListUsers(ctx, skip, limit)
}

// NewUserStorageLogging ...
func NewUserStorageLogging(

	logger *logrus.Logger,

	next user.Storage,
) user.Storage {
	return &logging{
		logger: logger,
		next:   next,
	}
}
