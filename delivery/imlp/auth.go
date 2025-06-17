package imlp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/delivery/dto"
	"github.com/johnquangdev/oauth2/middleware"
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
	m        middleware.MiddlewareCustom
}

func NewRegisterRouter(u interfaces.UseCase, g *echo.Group, v *validator.Validate, cfg utils.Config, m middleware.MiddlewareCustom) {
	r := implAuth{
		useCase:  u,
		validate: v,
		config:   cfg,
		m:        m,
	}
	google := g.Group("/google")

	google.GET("/callback", r.handlerCallback)
	google.GET("/login", r.handlerGetAuthURL)
	google.POST("/getme", r.handlerGetUserInfoFromGoogle)

	google.GET("/home", r.handlerMain)
	g.POST("/register", r.handlerCreateUser)
	g.POST("/logout", r.handlerLogout)
	g.GET("/profile", r.handleGetProfile, m.JWTAuthMiddleware())
}
func (implAuth implAuth) handlerMain(c echo.Context) error {
	htmlIndex := `<html>
				<body>
					<h2>OAuth2 Google Login với Echo Framework</h2>
					<a href="/v1/auth/google/login">Đăng nhập với Google</a>
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
		ID:         uuid.New(),
		Email:      user.Gmail,
		Name:       user.Name,
		Avatar:     user.Picture,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
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
	accessTokenTimeLife := time.Duration(u.config.AccessTokenTimeLife) * time.Hour
	refreshTokenTimeLife := time.Duration(u.config.RefreshTokenTimeLife) * time.Hour
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
	// Create new user
	if err := authUseCase.CreateUser(context.Background(), rModels.User{
		ID:         info.Id,
		Email:      info.Email,
		Name:       info.Name,
		Status:     "active",
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
		IsBlocked:             false,
		RefreshTokenExpiresAt: verifyRefreshToken.ExpiresAt.Time,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  verifyAccessToken.ExpiresAt.Time,
	}); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  http.StatusInternalServerError,
			"message": err.Error(),
		})
	}
	tokenRepose := dto.JwtResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		ExpiresAtAccessToken:  verifyAccessToken.ExpiresAt.Time,
		ExpiresAtRefreshToken: verifyRefreshToken.ExpiresAt.Time,
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":     http.StatusOK,
		"data user":  info,
		"data token": tokenRepose,
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
	userId := tokenVerify.Id
	// add user backlist
	if err := usecase.AddBackList(userId, logout.RefreshToken, time.Until((tokenVerify.ExpiresAt.Time))); err != nil {
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

func (u implAuth) handleGetProfile(c echo.Context) error {
	usecase, err := u.useCase.Auth()
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"status": http.StatusInternalServerError,
			"detail": err.Error(),
		})
	}
	userID, ok := c.Get("claims").(uuid.UUID)
	if !ok {
		fmt.Println(userID)
		return echo.NewHTTPError(http.StatusUnauthorized, "userID not found in context")
	}

	// authHeader := c.Request().Header.Get("Authorization")
	// if !strings.HasPrefix(authHeader, "Bearer ") {
	// 	return echo.NewHTTPError(http.StatusUnauthorized, "Missing or invalid Authorization header")
	// }

	// tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// claim, err := utils.VerifyToken(tokenString, u.config.SecretKey)
	// if err != nil {
	// 	return c.JSON(500, map[string]interface{}{
	// 		"status":  http.StatusInternalServerError,
	// 		"detail":  err.Error(),
	// 		"message": "token invalid",
	// 	})
	// }
	profile, err := usecase.GetUserByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(500, map[string]interface{}{
			"status": http.StatusInternalServerError,
			"detail": err.Error(),
		})
	}
	return c.JSON(200, map[string]interface{}{
		"status": http.StatusOK,
		"detail": "get ok",
		"data":   profile,
	})
}
