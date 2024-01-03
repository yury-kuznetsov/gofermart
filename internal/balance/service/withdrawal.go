package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"github.com/yury-kuznetsov/gofermart/validation"
	"time"
)

var ErrIncorrectOrder = errors.New("некорректный номер заказа")
var ErrInsufficientFunds = errors.New("на счету недостаточно средств")

type WithdrawalsRepository interface {
	Create(ctx context.Context, withdrawal model.Withdrawal) error
	FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error)
}

type WithdrawalService struct {
	bRepo BalanceRepository
	wRepo WithdrawalsRepository
}

func NewWithdrawalService(
	bRepo BalanceRepository,
	wRepo WithdrawalsRepository,
) *WithdrawalService {
	return &WithdrawalService{
		bRepo: bRepo,
		wRepo: wRepo,
	}
}

func (s *WithdrawalService) Withdraw(
	ctx context.Context,
	userID uuid.UUID,
	order string,
	sum float64,
) error {
	/*
		Нужен механизм транзакции, но описывать его тут кажется неуместно.
		Сервис не должен знать про всякие `db.Begin()` или `tx.Rollback()`.
		Где в таком случае его размещать - открытый вопрос к ревьюеру.

		tx, err := s.db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_ = tx.Rollback()
				return
			}
			err = tx.Commit()
		}()
	*/

	// проверяем корректность номера заказа
	if !validation.IsValidLuhn(order) {
		return ErrIncorrectOrder
	}

	// получаем баланс пользователя
	balance, err := s.bRepo.FindByUser(ctx, userID)
	if err != nil {
		return err
	}

	// проверяем наличие суммы для списания
	if (balance.Accrual - balance.Withdrawal) < sum {
		return ErrInsufficientFunds
	}

	// увеличиваем показатель использованных средств
	balance.Withdrawal += sum
	err = s.bRepo.Save(ctx, balance)
	if err != nil {
		return err
	}

	// сохраняем событие списания
	withdrawal := model.Withdrawal{
		ID:        uuid.New(),
		UserID:    userID,
		Number:    order,
		Sum:       sum,
		CreatedAt: time.Now(),
	}

	return s.wRepo.Create(ctx, withdrawal)
}

func (s *WithdrawalService) GetWithdrawals(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	return s.wRepo.FindByUser(ctx, userID)
}
