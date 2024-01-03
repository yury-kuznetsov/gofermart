package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/yury-kuznetsov/gofermart/internal/balance/mock"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	userID1 := uuid.New()
	userID2 := uuid.New()

	repo := &mock.AccrualRepo{}
	_ = repo.Save(context.Background(), model.Accrual{
		ID:        uuid.New(),
		UserID:    userID1,
		Number:    "12345678903",
		Status:    model.StatusNew,
		CreatedAt: time.Now(),
	})
	srv := &AccrualService{r: repo}

	tests := []struct {
		name   string
		userID uuid.UUID
		number string
		error  error
	}{
		{
			name:   "InvalidLuhn",
			userID: userID1,
			number: "123456789",
			error:  ErrIncorrectNumber,
		},
		{
			name:   "AlreadyLoadedByThisUser",
			userID: userID1,
			number: "12345678903",
			error:  ErrAlreadyLoadedByThisUser,
		},
		{
			name:   "AlreadyLoadedByThisUser",
			userID: userID2,
			number: "12345678903",
			error:  ErrAlreadyLoadedByAnotherUser,
		},
		{
			name:   "SuccessLoad",
			userID: userID1,
			number: "9278923470",
			error:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := srv.Load(context.Background(), tt.userID, tt.number)
			assert.Equal(t, tt.error, err)
		})
	}
}

func TestGetOrders(t *testing.T) {
	accrual := model.Accrual{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Number:    "12345678903",
		Status:    model.StatusNew,
		CreatedAt: time.Now(),
	}

	repo := &mock.AccrualRepo{}
	_ = repo.Save(context.Background(), accrual)
	srv := &AccrualService{r: repo}

	accruals, err := srv.GetOrders(context.Background(), accrual.UserID)
	assert.NoError(t, err)
	assert.Len(t, accruals, 1)
	assert.Equal(t, accruals[0].ID, accrual.ID)
}
