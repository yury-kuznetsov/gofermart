package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"strings"
)

type JWTService interface {
	GetUserID(token string) uuid.UUID
}

type key int

const (
	CookieKey     = "token"
	keyUserID key = iota
)

func AuthMiddleware(jwtService JWTService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// извлекаем токен из заголовка или куки
			tokenString := findToken(r)
			if tokenString == "" {
				http.Error(w, "Authorization required", http.StatusUnauthorized)
				return
			}

			// извлекаем идентификатор пользователя
			userID := jwtService.GetUserID(tokenString)
			if userID == uuid.Nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// передаем в контекст для обработчиков
			ctx := context.WithValue(r.Context(), keyUserID, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func findToken(r *http.Request) string {
	// ищем в заголовке
	tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	if tokenString != "" {
		return tokenString
	}

	// ищем в куки
	cookie, err := r.Cookie(CookieKey)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func GetUserID(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(keyUserID).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return id
}
