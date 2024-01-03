package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
)

type BalanceRepository interface {
	FindByUser(ctx context.Context, userID uuid.UUID) (model.Balance, error)
	Save(ctx context.Context, balance model.Balance) error
}

type BalanceService struct {
	r BalanceRepository
}

func NewBalanceService(bRepo BalanceRepository) *BalanceService {
	return &BalanceService{r: bRepo}
}

func (s *BalanceService) GetBalance(
	ctx context.Context,
	userID uuid.UUID,
) (model.Balance, error) {
	balance, err := s.r.FindByUser(ctx, userID)
	if err != nil {
		return model.Balance{}, err
	}

	// на клиенте ожидают текущий баланс пользователя
	balance.Accrual -= balance.Withdrawal

	return balance, nil
}
