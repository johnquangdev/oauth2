package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/middleware"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/usecase/models"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
)

type oAuth2GithubHandler struct {
	useCase    interfaces.UseCaseImpl
	validate   *validator.Validate
	config     utils.Config
	middleware middleware.MiddlewareCustom
}

func RegisterOAAuth2GithubHandler(u interfaces.UseCaseImpl, g *echo.Group, v *validator.Validate, cfg utils.Config, m middleware.MiddlewareCustom) {
	r := oAuth2GithubHandler{
		useCase:    u,
		validate:   v,
		config:     cfg,
		middleware: m,
	}
	github := g.Group("/github")

	github.GET("/login", r.handlerGithubLogin)
	github.GET("/callback", r.handlerGithubCallback)
}

func (h *oAuth2GithubHandler) handlerGithubLogin(c echo.Context) error {
	//call usecase to get login url
	loginURL, err := h.useCase.Auth().GithubOauth2.GetAuthURL()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (h *oAuth2GithubHandler) handlerGithubCallback(c echo.Context) error {
	var (
		ctx   = c.Request().Context()
		user  *models.User
		token *models.TokenJwt
	)
	//bind query param
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "Missing authorization code",
		})
	}
	//call usecase to login with github
	token, user, err := h.useCase.Auth().GithubOauth2.Login(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	//return token to client
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": http.StatusOK,
		"data": map[string]interface{}{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
			"user":          user,
		},
	})

}
