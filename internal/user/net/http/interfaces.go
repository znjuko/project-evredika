package http

import (
	"context"

	"github.com/labstack/echo"

	v1 "project-evredika/pkg/api/v1"
)

type logger interface {
	Error(args ...interface{})
}

type usecase interface {
	CreateUser(ctx context.Context, user *v1.User) (err error)
	DeleteUser(ctx context.Context, ID string) (err error)
	UpdateUser(ctx context.Context, user *v1.User) (err error)
	GetUser(ctx context.Context, ID string) (user *v1.User, err error)
	ListUsers(ctx context.Context, skip, limit int) (users []*v1.User)
}

type queryGetter interface {
	GetInt(ctx echo.Context, name string) (value int)
}
