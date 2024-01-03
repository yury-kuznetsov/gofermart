package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

const (
	SecretKey = "SECRET_KEY"
	Duration  = time.Hour
)

type JWTService struct{}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

func NewTokenService() *JWTService {
	return &JWTService{}
}

func (s *JWTService) GenerateToken(userID uuid.UUID) string {
	if userID == uuid.Nil {
		return ""
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(Duration)),
		},
		UserID: userID,
	})

	tokenString, _ := token.SignedString([]byte(SecretKey))

	return tokenString
}

func (s *JWTService) GetUserID(tokenString string) uuid.UUID {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(SecretKey), nil
	})

	if err != nil {
		return uuid.Nil
	}

	if !token.Valid {
		return uuid.Nil
	}

	return claims.UserID
}
