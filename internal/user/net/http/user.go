package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"

	"project-evredika/internal/user"
	v1 "project-evredika/pkg/api/v1"
)

const (
	idQuery    = "id"
	skipQuery  = "skip"
	limitQuery = "limit"

	urlCreateUser = "/api/v1/user"
	urlUpdateUser = "/api/v1/user"
	urlDeleteUser = "/api/v1/user"
	urlGetUser    = "/api/v1/user"
	urlListUsers  = "/api/v1/user/list"
)

type userDelivery struct {
	usecase     usecase
	queryGetter queryGetter
}

func (d *userDelivery) CreateUser(ctx echo.Context) (err error) {
	u := new(v1.User)

	var b []byte
	if b, err = ioutil.ReadAll(ctx.Request().Body); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()

	if err = json.Unmarshal(b, u); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err = d.usecase.CreateUser(ctx.Request().Context(), u); err != nil {
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}

func (d *userDelivery) UpdateUser(ctx echo.Context) (err error) {
	u := new(v1.User)

	var b []byte
	if b, err = ioutil.ReadAll(ctx.Request().Body); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}
	defer ctx.Request().Body.Close()

	if err = json.Unmarshal(b, u); err != nil {
		return ctx.NoContent(http.StatusBadRequest)
	}

	if err = d.usecase.UpdateUser(ctx.Request().Context(), u); err != nil {
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}

func (d *userDelivery) DeleteUser(ctx echo.Context) (err error) {
	if err = d.usecase.DeleteUser(ctx.Request().Context(), ctx.QueryParam(idQuery)); err != nil {
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}

func (d *userDelivery) GetUser(ctx echo.Context) (err error) {
	var u *v1.User
	if u, err = d.usecase.GetUser(ctx.Request().Context(), ctx.QueryParam(idQuery)); err != nil {
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK, u)
}

func (d *userDelivery) ListUsers(ctx echo.Context) (err error) {
	return ctx.JSON(http.StatusOK, d.usecase.ListUsers(
		ctx.Request().Context(),
		d.queryGetter.GetInt(ctx, skipQuery),
		d.queryGetter.GetInt(ctx, limitQuery),
	))
}

func (d *userDelivery) Initiate(server *echo.Echo) {
	server.POST(urlCreateUser, d.CreateUser)
	server.PUT(urlUpdateUser, d.UpdateUser)
	server.DELETE(urlDeleteUser, d.DeleteUser)
	server.GET(urlGetUser, d.GetUser)
	server.GET(urlListUsers, d.ListUsers)
}

// NewUserHttpDelivery ...
func NewUserHttpDelivery(
	usecase usecase,
	queryGetter queryGetter,
) user.Delivery {
	return &userDelivery{
		usecase:     usecase,
		queryGetter: queryGetter,
	}
}
