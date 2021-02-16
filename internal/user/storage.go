package user

import (
	"context"

	v1 "project-evredika/pkg/api/v1"
)

type Storage interface {
	CreateUser(ctx context.Context, user *v1.User) (err error)
	DeleteUser(ctx context.Context, ID string) (err error)
	UpdateUser(ctx context.Context, user *v1.User) (err error)
	GetUser(ctx context.Context, ID string) (user *v1.User, err error)
	ListUsers(ctx context.Context, skip, limit int) (users []*v1.User)
}
