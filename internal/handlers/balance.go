package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"github.com/yury-kuznetsov/gofermart/internal/balance/service"
	"github.com/yury-kuznetsov/gofermart/middleware"
	"io"
	"net/http"
)

type BalanceService interface {
	GetBalance(ctx context.Context, userID uuid.UUID) (model.Balance, error)
}

type AccrualService interface {
	Load(ctx context.Context, userID uuid.UUID, number string) error
	GetOrders(ctx context.Context, userID uuid.UUID) ([]model.Accrual, error)
}

type WithdrawalService interface {
	Withdraw(ctx context.Context, userID uuid.UUID, order string, sum float64) error
	GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error)
}

type withdrawRequest struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

func GetBalanceHandler(s BalanceService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())

		// получаем баланс пользователя
		balance, err := s.GetBalance(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// возвращаем ответ
		w.Header().Set("content-type", "application/json")
		err = json.NewEncoder(w).Encode(balance)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func LoadNumberHandler(s AccrualService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())

		// получение номера заказа
		number, err := io.ReadAll(r.Body)
		if err != nil || len(number) == 0 {
			http.Error(w, "неверный формат запроса", http.StatusBadRequest)
			return
		}

		// загрузка номера заказа
		err = s.Load(r.Context(), userID, string(number))
		if err != nil {
			switch {
			case errors.Is(err, service.ErrAlreadyLoadedByThisUser):
				http.Error(w, err.Error(), http.StatusOK)
			case errors.Is(err, service.ErrAlreadyLoadedByAnotherUser):
				http.Error(w, err.Error(), http.StatusConflict)
			case errors.Is(err, service.ErrIncorrectNumber):
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func GetOrdersHandler(s AccrualService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())

		// получение заказов пользователя
		orders, err := s.GetOrders(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// проверка наличия заказов
		if len(orders) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// преобразование заказов в JSON
		w.Header().Set("content-type", "application/json")
		if err = json.NewEncoder(w).Encode(orders); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func WithdrawHandler(s WithdrawalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())

		// принимаем запрос
		var request withdrawRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "не переданы order или sum", http.StatusBadRequest)
			return
		}

		err = s.Withdraw(r.Context(), userID, request.Order, request.Sum)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrIncorrectOrder):
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			case errors.Is(err, service.ErrInsufficientFunds):
				http.Error(w, err.Error(), http.StatusPaymentRequired)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func GetWithdrawalsHandler(s WithdrawalService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := middleware.GetUserID(r.Context())

		// получение списаний пользователя
		withdrawals, err := s.GetWithdrawals(r.Context(), userID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// проверка наличия записей
		if len(withdrawals) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// преобразование заказов в JSON
		w.Header().Set("content-type", "application/json")
		if err := json.NewEncoder(w).Encode(withdrawals); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
