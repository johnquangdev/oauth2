package imlp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/delivery/dto"
	rModels "github.com/johnquangdev/oauth2/repository/models"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/usecase/models"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
)

type implAuth struct {
	validate *validator.Validate
	useCase  interfaces.UseCase
	config   utils.Config
}

func NewRegisterRouter(u interfaces.UseCase, g *echo.Group, v *validator.Validate, cfg utils.Config) {
	r := implAuth{
		useCase:  u,
		validate: v,
		config:   cfg,
	}
	g.GET("/callback", r.handlerCallback)
	g.GET("/login", r.handlerGetAuthURL)
	g.GET("/home", r.handlerMain)
	g.POST("/register", r.handlerCreateUser)
	g.POST("/getuser", r.handlerGetUserInfoFromGoogle)
	g.POST("/logout", r.handlerLogout)
}
func (implAuth implAuth) handlerMain(c echo.Context) error {
	htmlIndex := `<html>
<body>
	<h2>OAuth2 Google Login với Echo Framework</h2>
	<a href="/v1/auth/login">Đăng nhập với Google</a>
</body>
</html>`
	return c.HTML(http.StatusOK, htmlIndex)
}

func (u *implAuth) handlerGetAuthURL(c echo.Context) error {
	authUseCase, err := u.useCase.Auth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to get auth usecase")
	}
	state, err := authUseCase.GenerateRandomState()
	if err != nil {
		return err
	}
	loginURL, err := authUseCase.GetGoogleAuthURL(state)
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusTemporaryRedirect, loginURL)
}

func (u *implAuth) handlerCallback(c echo.Context) error {
	authUseCase, err := u.useCase.Auth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to get auth usecase")
	}
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

	response, err := authUseCase.ExchangeCodeForToken(c.Request().Context(), req)
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

func (u *implAuth) handlerCreateUser(c echo.Context) error {
	authUseCase, err := u.useCase.Auth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to get auth usecase")
	}
	var user dto.User
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err,
		})
	}
	if err := u.validate.Struct(&user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err,
		})
	}
	if err := authUseCase.CreateUser(context.Background(), rModels.User{
		ID:     uuid.New(),
		Email:  user.Gmail,
		Name:   user.Name,
		Avatar: user.Picture,
	}); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err,
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "create user ok",
	})
}

func (u *implAuth) handlerGetUserInfoFromGoogle(c echo.Context) error {
	// get the method from layer usecase
	authUseCase, err := u.useCase.Auth()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "failed to get auth usecase")
	}
	// bind data from FE
	var userInfo dto.GetUserInfo
	if err := c.Bind(&userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err,
		})
	}
	// check validate data from fe
	if err := u.validate.Struct(&userInfo); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// use the accesstoken to get userinfo from google api
	info, err := authUseCase.GetUserInfoGoogle(userInfo.AccessToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"status":  http.StatusUnauthorized,
			"message": err.Error(),
		})
	}
	// check user exited, if user exited return token and profile
	// else if create new user with token and return
	// err = authUseCase.IsUserExists(info.Gmail)
	// if err != nil {
	// 	return c.JSON(http.StatusBadRequest, map[string]interface{}{
	// 		"status":  http.StatusBadRequest,
	// 		"message": err.Error(),
	// 	})
	// }
	accessTokenTimeLife := time.Duration(u.config.AccessTokenTimeLife) * time.Hour
	refreshTokenTimeLife := time.Duration(u.config.RefreshTokenTimeLife) * time.Hour
	fmt.Println(accessTokenTimeLife, refreshTokenTimeLife)
	// ganerate token jwt
	accessToken := utils.GenerateToken(info.Id, info.Email, info.Name, accessTokenTimeLife, u.config.SecretKey)
	refreshToken := utils.GenerateToken(info.Id, info.Email, info.Name, refreshTokenTimeLife, u.config.SecretKey)
	// verify token jwt
	verifyAccessToken, err := utils.VerifyToken(accessToken, u.config.SecretKey)
	if err != nil {
		return fmt.Errorf("verify access token token err: %v", err)
	}
	verifyRefreshToken, err := utils.VerifyToken(refreshToken, u.config.SecretKey)
	if err != nil {
		return fmt.Errorf("verify refresh token token err: %v", err)
	}
	if err := authUseCase.SaveRefreshToken(info.ProviderID, refreshToken, time.Minute*30); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// Create new user
	if err := authUseCase.CreateUser(context.Background(), rModels.User{
		ID:         info.Id,
		Email:      info.Email,
		Name:       info.Name,
		Avatar:     info.Avatar,
		ProviderID: info.ProviderID,
		Provider:   info.Provider,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	//Create session
	if err := authUseCase.CreateSesion(context.Background(), rModels.Session{
		ID:                    uuid.NewString(),
		UserID:                info.Id.String(),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: verifyRefreshToken.ExpiresAt.Time,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  verifyAccessToken.ExpiresAt.Time,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":                http.StatusOK,
		"data":                  info,
		"accessToken":           accessToken,
		"refreshToken":          refreshToken,
		"AccessTokenExpiresAt":  verifyAccessToken.ExpiresAt.Time,
		"RefreshTokenExpiresAt": verifyRefreshToken.ExpiresAt.Time,
	})
}
func (u *implAuth) handlerLogout(c echo.Context) error {
	// call logic useCase
	usecase, err := u.useCase.Auth()
	if err != nil {
		return err
	}
	var logout dto.Logout
	// bind refreshToken
	if err := c.Bind(&logout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// check validate
	if err := u.validate.Struct(&logout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// verify refresh token
	tokenVerify, err := utils.VerifyToken(logout.RefreshToken, u.config.SecretKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// add user backlist
	if err := usecase.SaveRefreshToken(tokenVerify.ID, logout.RefreshToken, time.Until((tokenVerify.ExpiresAt.Time))); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// logic logout
	if err := usecase.Logout(context.Background(), tokenVerify.Id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// gọi logic usecase
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "logout ok",
	})
}
