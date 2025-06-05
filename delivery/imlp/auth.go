package imlp

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/usecase/models"
	"github.com/labstack/echo/v4"
)

type implAuth struct {
	validate *validator.Validate
	usecase  interfaces.UseCase
}

func NewRegisterRouter(u interfaces.UseCase, g *echo.Group, v *validator.Validate) {
	r := implAuth{
		usecase:  u,
		validate: v,
	}
	g.GET("/callback", r.HandleCallback)
	g.GET("/login", r.GetAuthURL)
	g.GET("/home", r.handleMain)
}
func (implAuth implAuth) handleMain(c echo.Context) error {
	htmlIndex := `<html>
<body>
	<h2>OAuth2 Google Login với Echo Framework</h2>
	<a href="/v1/auth/login">Đăng nhập với Google</a>
</body>
</html>`
	return c.HTML(http.StatusOK, htmlIndex)
}

func (h *implAuth) GetAuthURL(c echo.Context) error {
	state, err := h.usecase.Auth().GenerateRandomState()
	if err != nil {
		return err
	}
	loginURL, err := h.usecase.Auth().GetGoogleAuthURL(state)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (u *implAuth) HandleCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")

	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "authorization code not found",
		})
	}

	req := models.ExchangeTokenRequest{
		Code:  code,
		State: state,
	}

	response, err := u.usecase.Auth().ExchangeCodeForToken(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login thành công!",
		"data":    response,
	})
}
