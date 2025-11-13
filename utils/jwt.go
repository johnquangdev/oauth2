package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type myCustomClaim struct {
	Id    uuid.UUID
	Name  string
	Gmail string
	jwt.RegisteredClaims
}

// GenerateToken tạo JWT với thời gian hết hạn UTC và trả về expires_in (giây)
func GenerateToken(id uuid.UUID, name string, gmail string, tokenTimeLife time.Duration, cfg string) (string, myCustomClaim) {
	expiresAt := time.Now().UTC().Add(tokenTimeLife)
	claims := myCustomClaim{
		Id:    id,
		Name:  name,
		Gmail: gmail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(cfg))
	if err != nil {
		return err.Error(), myCustomClaim{}
	}
	return t, claims
}

func VerifyToken(tokenStr string, secretKey string) (*myCustomClaim, error) {
	if strings.TrimSpace(tokenStr) == "" {
		return nil, fmt.Errorf("token is empty")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &myCustomClaim{}, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token by err: %v", err)
	}

	claims, ok := token.Claims.(*myCustomClaim)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token by err: %v", err)
	}
	// Kiểm tra thời hạn
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token expiresAt")
	}

	return claims, nil
}
