package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yury-kuznetsov/gofermart/internal/balance/mock"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"testing"
)

func TestGetBalance(t *testing.T) {
	row := model.Balance{
		UserID:     uuid.New(),
		Accrual:    100,
		Withdrawal: 20,
	}

	repo := &mock.BalanceRepo{}
	_ = repo.Save(context.Background(), row)
	srv := &BalanceService{r: repo}

	balance, err := srv.GetBalance(context.Background(), row.UserID)
	assert.NoError(t, err)
	assert.Equal(t, balance.Accrual, row.Accrual-row.Withdrawal)
	assert.Equal(t, balance.Withdrawal, row.Withdrawal)
}
