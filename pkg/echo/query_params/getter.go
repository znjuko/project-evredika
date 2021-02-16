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
	minimal    int
}

func (g *getter) GetInt(ctx echo.Context, name string) (value int) {
	value, err := strconv.Atoi(ctx.QueryParam(name))
	if err != nil {
		return g.intDefault
	}

	if value < g.minimal {
		return g.minimal
	}

	return value
}

// NewQueryGetter ...
func NewQueryGetter(intDefault, minimal int) Getter {
	return &getter{
		intDefault: intDefault,
		minimal:    minimal,
	}
}
