package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"time"
)

type WithdrawalRepository struct {
	db *sql.DB
}

func NewWithdrawalRepository(db *sql.DB) *WithdrawalRepository {
	r := &WithdrawalRepository{db: db}

	_, err := r.db.Exec(`CREATE TABLE IF NOT EXISTS balance_withdrawal (
		id         uuid      not null constraint withdrawals_pk primary key,
		user_id    uuid      not null constraint withdrawals_users_id_fk references "user",
		number     varchar   not null constraint withdrawals_pk_2 unique,
		sum        decimal   not null,
		created_at timestamp
	)`)

	if err != nil {
		fmt.Println(err.Error())
	}

	return r
}

func (r *WithdrawalRepository) Create(ctx context.Context, model model.Withdrawal) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO balance_withdrawal VALUES ($1, $2, $3, $4, $5)`,
		model.ID, model.UserID, model.Number, model.Sum, model.CreatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *WithdrawalRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Withdrawal, error) {
	var withdrawals []model.Withdrawal
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT * FROM balance_withdrawal WHERE user_id = $1 ORDER BY created_at",
		userID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var withdrawal model.Withdrawal
		err := rows.Scan(
			&withdrawal.ID,
			&withdrawal.UserID,
			&withdrawal.Number,
			&withdrawal.Sum,
			&withdrawal.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}
