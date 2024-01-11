package mock

import (
	"context"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
)

type WithdrawalRepo struct {
	withdrawals []model.Withdrawal
}

func (w *WithdrawalRepo) Create(_ context.Context, withdrawal model.Withdrawal) error {
	w.withdrawals = append(w.withdrawals, withdrawal)
	return nil
}

func (w *WithdrawalRepo) FindByUser(_ context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal
	for _, withdrawal := range w.withdrawals {
		if withdrawal.UserID == userID {
			withdrawals = append(withdrawals, withdrawal)
		}
	}
	return withdrawals, nil
}
