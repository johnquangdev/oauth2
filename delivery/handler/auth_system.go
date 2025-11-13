package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/johnquangdev/oauth2/delivery/models"
	"github.com/johnquangdev/oauth2/middleware"
	"github.com/johnquangdev/oauth2/usecase/interfaces"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
)

type AuthSystemHandler struct {
	validate   *validator.Validate
	useCase    interfaces.UseCaseImpl
	config     utils.Config
	middleware middleware.MiddlewareCustom
}

func RegisterAuthSystemHandler(u interfaces.UseCaseImpl, g *echo.Group, v *validator.Validate, cfg utils.Config, m middleware.MiddlewareCustom) {
	r := &AuthSystemHandler{
		useCase:    u,
		validate:   v,
		config:     cfg,
		middleware: m,
	}
	g.POST("/logout", r.handlerLogout)
	g.GET("/profile", r.handleGetProfile, m.JWTAuthMiddleware())
}

// GetUserProfile godoc
// @Summary Lấy thông tin người dùng
// @Description Trả về thông tin người dùng đã xác thực qua Bearer Token (JWT)
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {object} rModels.User
// @Failure 401 {object} map[string]interface{} "Unauthorized – Token không hợp lệ hoặc hết hạn"
// @Failure 500 {object} map[string]interface{} "Lỗi hệ thống"
// @Router /v1/auth/profile [get]
func (h *AuthSystemHandler) handleGetProfile(c echo.Context) error {
	userId, ok := c.Get("claims").(uuid.UUID)
	if !ok {
		fmt.Println(userId)
		return echo.NewHTTPError(http.StatusUnauthorized, "userId not found in context")
	}

	profile, err := h.useCase.Auth().SystemAuth.GetUserById(context.Background(), userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status": http.StatusInternalServerError,
			"detail": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, profile)
}

// @Summary Logout người dùng
// @Description Thoát phiên làm việc của người dùng
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /v1/auth/logout [post]
func (h *AuthSystemHandler) handlerLogout(c echo.Context) error {
	var logout models.Logout
	// bind refreshToken
	if err := c.Bind(&logout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// check validate
	if err := h.validate.Struct(&logout); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	// verify refresh token
	tokenVerify, err := utils.VerifyToken(logout.RefreshToken, h.config.SecretKey)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}
	userId := tokenVerify.Id
	// call usecase logout
	if err := h.useCase.Auth().GoogleOauth2.Logout(context.Background(), tokenVerify.Id); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	// add user backlist
	if err := h.useCase.Auth().SystemAuth.AddBackList(userId, logout.RefreshToken, time.Until((tokenVerify.ExpiresAt.Time))); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"status":  http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  http.StatusOK,
		"message": "logout ok",
	})
}
