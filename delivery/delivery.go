package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/delivery/imlp"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/labstack/echo/v4"
)

func NewDelivery(u interfaces.UseCase, g *echo.Group, v *validator.Validate) error {
	a := g.Group("/auth")
	imlp.NewRegisterRouter(u, a, v)
	return nil
}
