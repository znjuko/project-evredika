package query_params

import (
	"strconv"

	"github.com/labstack/echo"
)

type Getter interface {
	GetInt(ctx echo.Context, name string) (value int)
}

type getter struct {
	intDefault int
}

func (g *getter) GetInt(ctx echo.Context, name string) (value int) {
	value, err := strconv.Atoi(ctx.QueryParam(name))
	if err != nil {
		return g.intDefault
	}

	return value
}

// NewQueryGetter ...
func NewQueryGetter(intDefault int) Getter { return &getter{intDefault: intDefault} }
