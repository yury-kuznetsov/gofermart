package model

import "github.com/google/uuid"

type Balance struct {
	UserID     uuid.UUID `json:"-"`
	Accrual    float64   `json:"current"`
	Withdrawal float64   `json:"withdrawn"`
}
