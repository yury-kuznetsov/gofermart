package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/user/service"
	"github.com/yury-kuznetsov/gofermart/middleware"
	"net/http"
)

type UserService interface {
	Register(ctx context.Context, login, password string) (uuid.UUID, error)
	Login(ctx context.Context, login, password string) (uuid.UUID, error)
}

type JWTService interface {
	GenerateToken(userID uuid.UUID) string
}

type registerRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterHandler(userService UserService, jwtService JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// принимаем запрос
		var request registerRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "не переданы login или password", http.StatusBadRequest)
			return
		}

		// регистрируем пользователя
		userID, err := userService.Register(r.Context(), request.Login, request.Password)
		if err != nil {
			if errors.Is(err, service.ErrUserExists) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// генерируем токен и сохраняем в куки
		token := jwtService.GenerateToken(userID)
		http.SetCookie(w, &http.Cookie{
			Name:  middleware.CookieKey,
			Value: token,
		})
		w.Header().Set("Authorization", token)

		w.WriteHeader(http.StatusOK)
	}
}

func LoginHandler(userService UserService, jwtService JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// принимаем запрос
		var request loginRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "не переданы login или password", http.StatusBadRequest)
			return
		}

		// авторизуем пользователя
		userID, err := userService.Login(r.Context(), request.Login, request.Password)
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// генерируем токен и сохраняем в куки
		token := jwtService.GenerateToken(userID)
		http.SetCookie(w, &http.Cookie{
			Name:  middleware.CookieKey,
			Value: token,
		})
		w.Header().Set("Authorization", token)

		w.WriteHeader(http.StatusOK)
	}
}
