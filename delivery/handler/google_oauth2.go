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

type oAuth2GoogleHandler struct {
	validate   *validator.Validate
	useCase    interfaces.UseCaseImpl
	config     utils.Config
	middleware middleware.MiddlewareCustom
}

func RegisterOAuth2GoogleHandler(u interfaces.UseCaseImpl, g *echo.Group, v *validator.Validate, cfg utils.Config, m middleware.MiddlewareCustom) {
	r := oAuth2GoogleHandler{
		useCase:    u,
		validate:   v,
		config:     cfg,
		middleware: m,
	}
	google := g.Group("/google")

	google.GET("/login", r.handlerGoogleLogin)
	google.GET("/callback", r.handlerGoogleCallback)

}

// @Summary Google Login
// @Description Trả về URL để redirect người dùng đến trang đăng nhập của Google
// @Tags OAuth2
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /v1/auth/google/login [get]
func (h *oAuth2GoogleHandler) handlerGoogleLogin(c echo.Context) error {
	loginURL, err := h.useCase.Auth().GoogleOauth2.GetAuthURL()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

// @Summary Google Callback
// @Description Callback and create user or allow user login
// @Tags OAuth2
// @Accept  json
// @Produce  json
// @Param token body map[string]string true "token từ Google"
// @Success 200 {object} dModel.JwtResponse
// @Failure 400 {object} map[string]interface{}
// @Router /v1/auth/google/callback [get]
func (h *oAuth2GoogleHandler) handlerGoogleCallback(c echo.Context) error {
	var (
		ctx   = c.Request().Context()
		user  *models.User
		token *models.TokenJwt
	)
	// bind data from FE
	code := c.QueryParam("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": "Missing authorization code",
		})
	}
	token, user, err := h.useCase.Auth().GoogleOauth2.Login(ctx, code)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": http.StatusInternalServerError,
			"detail": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":   http.StatusOK,
		"token":    token,
		"userinfo": user,
	})
}
