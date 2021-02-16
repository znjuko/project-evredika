package storage

import (
	"context"
	"time"

	"project-evredika/internal/user"
	v1 "project-evredika/pkg/api/v1"
)

type logging struct {
	logger logger

	next user.Storage
}

func (l *logging) CreateUser(ctx context.Context, user *v1.User) (err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "Create User",
			"request", user,
			"error", err,
		)
	}(time.Now())
	return l.next.CreateUser(ctx, user)
}

func (l *logging) DeleteUser(ctx context.Context, ID string) (err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "Delete User",
			"request", ID,
			"error", err,
		)
	}(time.Now())
	return l.next.DeleteUser(ctx, ID)
}

func (l *logging) UpdateUser(ctx context.Context, user *v1.User) (err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "Update User",
			"request", user,
			"error", err,
		)
	}(time.Now())
	return l.next.UpdateUser(ctx, user)
}

func (l *logging) GetUser(ctx context.Context, ID string) (user *v1.User, err error) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "Get User",
			"request", ID,
			"response", user,
			"error", err,
		)
	}(time.Now())
	return l.next.GetUser(ctx, ID)
}

func (l *logging) ListUsers(ctx context.Context, skip, limit int) (users []*v1.User) {
	defer func(begin time.Time) {
		l.logger.Debug(
			"method", "List User",
			"skip", skip,
			"limit", limit,
			"response", users,
		)
	}(time.Now())
	return l.next.ListUsers(ctx, skip, limit)
}

// NewUserStorageLogging ...
func NewUserStorageLogging(

	logger logger,

	next user.Storage,
) user.Storage {
	return &logging{
		logger: logger,
		next:   next,
	}
}
