package user

import (
	"github.com/labstack/echo"
)

type Delivery interface {
	CreateUser(ctx echo.Context) (err error)
	DeleteUser(ctx echo.Context) (err error)
	UpdateUser(ctx echo.Context) (err error)
	GetUser(ctx echo.Context) (err error)
	ListUsers(ctx echo.Context) (err error)
	// used for initiating echo-handlers
	Initiate(server *echo.Echo)
}
