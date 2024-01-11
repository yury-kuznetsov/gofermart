package mock

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
)

type BalanceRepo struct {
	balances []model.Balance
}

func (b *BalanceRepo) Save(_ context.Context, model model.Balance) error {
	for i, balance := range b.balances {
		if balance.UserID == model.UserID {
			b.balances[i] = model
			return nil
		}
	}
	b.balances = append(b.balances, model)
	return nil
}

func (b *BalanceRepo) FindByUser(_ context.Context, userID uuid.UUID) (model.Balance, error) {
	for _, balance := range b.balances {
		if balance.UserID == userID {
			return balance, nil
		}
	}
	return model.Balance{}, errors.New("balance not found")
}
