package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yury-kuznetsov/gofermart/internal/balance/mock"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"testing"
)

func TestWithdraw(t *testing.T) {
	userID := uuid.New()
	balance := model.Balance{
		UserID:     userID,
		Accrual:    100,
		Withdrawal: 20,
	}

	bRepo := &mock.BalanceRepo{}
	_ = bRepo.Save(context.Background(), balance)
	wRepo := &mock.WithdrawalRepo{}
	srv := &WithdrawalService{bRepo: bRepo, wRepo: wRepo}

	tests := []struct {
		name   string
		number string
		sum    float64
		error  error
	}{
		{
			name:   "InvalidLuhn",
			number: "123456789",
			sum:    10,
			error:  ErrIncorrectOrder,
		},
		{
			name:   "InsufficientFunds",
			number: "12345678903",
			sum:    90, // на счету 100-20=80 баллов
			error:  ErrInsufficientFunds,
		},
		{
			name:   "Success",
			number: "12345678903",
			sum:    70,
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.Withdraw(context.Background(), userID, tt.number, tt.sum)
			assert.Equal(t, tt.error, err)

			if err == nil {
				// после последнего успешного списания 70 баллов
				row, err := bRepo.FindByUser(context.Background(), userID)
				assert.NoError(t, err)
				assert.Equal(t, balance.Accrual, row.Accrual)
				assert.Equal(t, balance.Withdrawal+tt.sum, row.Withdrawal)
			}
		})
	}

}
