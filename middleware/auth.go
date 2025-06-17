package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/johnquangdev/oauth2/repository/interfaces"
	"github.com/johnquangdev/oauth2/utils"
	"github.com/labstack/echo/v4"
)

type MiddlewareCustom struct {
	cfg  utils.Config
	repo interfaces.Repo
}

func NewMiddleware(cfg utils.Config, repo interfaces.Repo) MiddlewareCustom {
	return MiddlewareCustom{
		cfg:  cfg,
		repo: repo,
	}
}

func (m MiddlewareCustom) JWTAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Lấy Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "Authorization header is required",
				})
			}

			// Kiểm tra format "Bearer <token>"
			if !strings.HasPrefix(authHeader, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "Invalid authorization header format. Use 'Bearer <token>'",
				})
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "Token is required",
				})
			}

			// Validate JWT token
			claims, err := utils.VerifyToken(tokenString, m.cfg.SecretKey)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"detail":  err.Error(),
					"message": "token invalid",
				})
			}

			// Kiểm tra token có bị blacklist không (Redis)
			isBlacklisted, err := m.repo.Redis().IsTokenBlacklisted(claims.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{
					"error": "error checking token status",
				})
			}
			if isBlacklisted {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "token is blocked",
				})
			}
			// // Lấy user info từ database
			user, err := m.repo.Auth().GetUserByUserId(context.Background(), claims.Id)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{
					"error": "user not found",
				})
			}

			// Kiểm tra user status
			if user.Status != "active" {
				return echo.NewHTTPError(http.StatusForbidden, map[string]string{
					"error": "account is suspended",
				})
			}

			// // Set user context để các handler khác sử dụng
			// c.Set("user", user)
			c.Set("claims", claims.Id)
			return next(c)
		}
	}
}
