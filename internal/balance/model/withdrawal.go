package model

import (
	"github.com/google/uuid"
	"time"
)

type Withdrawal struct {
	ID        uuid.UUID `json:"-"`
	UserID    uuid.UUID `json:"-"`
	Number    string    `json:"order"`
	Sum       float64   `json:"sum"`
	CreatedAt time.Time `json:"processed_at"`
}
