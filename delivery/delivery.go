package delivery

import (
	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/delivery/handler"
	"github.com/johnquangdev/oauth2/middleware"
	uInterface "github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
)

func NewDelivery(u uInterface.UseCaseImpl, g *echo.Group, v *validator.Validate, cfg utils.Config, m middleware.MiddlewareCustom) {
	auth := g.Group("/auth")
	handler.RegisterAuthSystemHandler(u, auth, v, cfg, m)
	handler.RegisterOAuth2GoogleHandler(u, auth, v, cfg, m)
	handler.RegisterOAAuth2GithubHandler(u, auth, v, cfg, m)
}
