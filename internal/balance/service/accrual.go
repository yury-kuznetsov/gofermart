package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/cmd/gophermart/config"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"github.com/yury-kuznetsov/gofermart/validation"
	"time"
)

var ErrIncorrectNumber = errors.New("некорректный номер заказа")
var ErrAlreadyLoadedByThisUser = errors.New("номер заказа уже был загружен этим пользователем")
var ErrAlreadyLoadedByAnotherUser = errors.New("номер заказа уже был загружен другим пользователем")

type AccrualRepository interface {
	Save(ctx context.Context, model model.Accrual) error
	FindByNumber(ctx context.Context, number string) (model.Accrual, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Accrual, error)
	FindForSync(ctx context.Context) ([]model.Accrual, error)
}

type AccrualService struct {
	r AccrualRepository
}

func NewAccrualService(bRepo BalanceRepository, aRepo AccrualRepository) *AccrualService {
	// запускаем сервис синхронизации
	go NewSyncService(bRepo, aRepo, config.Options.AccrualAddr).Start()

	return &AccrualService{r: aRepo}
}

func (s *AccrualService) Load(ctx context.Context, userID uuid.UUID, number string) error {
	// проверяем номер по алгоритму Луна
	if !validation.IsValidLuhn(number) {
		return ErrIncorrectNumber
	}

	// проверяем наличие заказа с таким номером
	accrual, _ := s.r.FindByNumber(ctx, number)
	if accrual.ID != uuid.Nil {
		if accrual.UserID == userID {
			return ErrAlreadyLoadedByThisUser
		}
		return ErrAlreadyLoadedByAnotherUser
	}

	// добавляем номер заказа
	accrual = model.Accrual{
		ID:        uuid.New(),
		UserID:    userID,
		Number:    number,
		Status:    model.StatusNew,
		Sum:       nil,
		CreatedAt: time.Now(),
	}
	err := s.r.Save(ctx, accrual)

	if err != nil {
		return err
	}

	return nil
}

func (s *AccrualService) GetOrders(ctx context.Context, userID uuid.UUID) ([]model.Accrual, error) {
	return s.r.FindByUser(ctx, userID)
}
