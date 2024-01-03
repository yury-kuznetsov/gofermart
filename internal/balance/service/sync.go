package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"net/http"
	"strconv"
	"time"
)

type SyncService interface {
	Start()
}

type errTooManyRequests struct {
	RetryAfter int
}

func (e *errTooManyRequests) Error() string {
	return fmt.Sprintf("Следующий запрос будет через %d секунд.", e.RetryAfter)
}

type syncService struct {
	bRepo BalanceRepository
	aRepo AccrualRepository
	host  string
}

func NewSyncService(
	bRepo BalanceRepository,
	aRepo AccrualRepository,
	host string,
) SyncService {
	return &syncService{
		bRepo: bRepo,
		aRepo: aRepo,
		host:  host,
	}
}

func (s *syncService) Start() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		<-ticker.C
		orders, err := s.aRepo.FindForSync(context.Background())
		if err != nil {
			fmt.Println(err)
			continue
		}

		for _, order := range orders {
			err := processOrder(s, order)
			if err != nil {
				fmt.Println(err)

				// после 429 ошибки ждем какое-то время
				var e *errTooManyRequests
				if errors.As(err, &e) {
					ticker.Stop()
					time.Sleep(time.Duration(e.RetryAfter) * time.Second)
					ticker = time.NewTicker(5 * time.Second)
				}
			}
		}
	}
}

func processOrder(s *syncService, order model.Accrual) error {
	resp, err := http.Get(s.host + "/api/orders/" + order.Number)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// превышено количество запросов к сервису
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter, _ := strconv.Atoi(resp.Header.Get("Retry-After"))
		return &errTooManyRequests{RetryAfter: retryAfter}
	}

	// заказ не зарегистрирован в системе расчёта
	if resp.StatusCode == http.StatusNoContent {
		order.Status = model.StatusInvalid
	}

	// внутренняя ошибка сервера
	if resp.StatusCode == http.StatusInternalServerError {
		order.Status = model.StatusInvalid
	}

	if resp.StatusCode == http.StatusOK {
		var respBody struct {
			Status  string   `json:"status"`
			Accrual *float64 `json:"accrual,omitempty"`
		}
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return err
		}

		if respBody.Status == "REGISTERED" {
			order.Status = model.StatusProcessing
		}

		if respBody.Status == "INVALID" {
			order.Status = model.StatusInvalid
		}

		if respBody.Status == "PROCESSING" {
			order.Status = model.StatusProcessing
		}

		if respBody.Status == "PROCESSED" {
			order.Status = model.StatusProcessed
			order.Sum = respBody.Accrual

			// меняем баланс пользователю
			err = changeBalance(s, order.UserID, *respBody.Accrual)
			if err != nil {
				return err
			}
		}
	}

	return s.aRepo.Save(context.Background(), order)
}

func changeBalance(s *syncService, userID uuid.UUID, accrual float64) error {
	balance, err := s.bRepo.FindByUser(context.Background(), userID)
	if err != nil {
		return err
	}

	balance.Accrual += accrual
	return s.bRepo.Save(context.Background(), balance)
}
