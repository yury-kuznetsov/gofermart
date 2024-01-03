package mock

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
)

type AccrualRepo struct {
	accruals []model.Accrual
}

func (a *AccrualRepo) Save(_ context.Context, model model.Accrual) error {
	for i, accrual := range a.accruals {
		if accrual.ID == model.ID {
			a.accruals[i] = model
			return nil
		}
	}
	a.accruals = append(a.accruals, model)
	return nil
}

func (a *AccrualRepo) FindByNumber(_ context.Context, number string) (model.Accrual, error) {
	for _, accrual := range a.accruals {
		if accrual.Number == number {
			return accrual, nil
		}
	}
	return model.Accrual{}, errors.New("accrual not found")
}

func (a *AccrualRepo) FindByUser(_ context.Context, userID uuid.UUID) ([]model.Accrual, error) {
	var accruals []model.Accrual
	for _, accrual := range a.accruals {
		if accrual.UserID == userID {
			accruals = append(accruals, accrual)
		}
	}
	return accruals, nil
}

func (a *AccrualRepo) FindForSync(_ context.Context) ([]model.Accrual, error) {
	var accruals []model.Accrual
	for _, accrual := range a.accruals {
		if accrual.Status == model.StatusNew {
			accruals = append(accruals, accrual)
		}
	}
	return accruals, nil
}
