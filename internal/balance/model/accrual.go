package model

import (
	"github.com/google/uuid"
	"time"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type Accrual struct {
	ID        uuid.UUID `json:"-"`
	UserID    uuid.UUID `json:"-"`
	Number    string    `json:"number"`
	Status    string    `json:"status"`
	Sum       *float64  `json:"accrual"`
	CreatedAt time.Time `json:"uploaded_at"`
}
