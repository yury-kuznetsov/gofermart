package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
)

type BalanceRepository struct {
	db *sql.DB
}

func NewBalanceRepository(db *sql.DB) *BalanceRepository {
	r := &BalanceRepository{db: db}

	_, _ = r.db.Exec(`CREATE TABLE IF NOT EXISTS balance (
		user_id    uuid    not null constraint balance_pk unique constraint balance_users_id_fk references "user",
		accrual    decimal not null,
		withdrawal decimal not null
	)`)

	return r
}

func (r *BalanceRepository) FindByUser(ctx context.Context, userID uuid.UUID) (model.Balance, error) {
	balance := model.Balance{UserID: userID, Accrual: 0, Withdrawal: 0}
	err := r.db.QueryRowContext(
		ctx,
		"SELECT accrual, withdrawal FROM balance WHERE user_id = $1",
		userID,
	).Scan(&balance.Accrual, &balance.Withdrawal)

	if errors.Is(err, sql.ErrNoRows) {
		return balance, nil
	}

	return balance, err
}

func (r *BalanceRepository) Save(ctx context.Context, balance model.Balance) error {
	// обновляем баланс пользователя
	result, err := r.db.ExecContext(
		ctx,
		"UPDATE balance SET accrual = $1, withdrawal = $2 WHERE user_id = $3",
		balance.Accrual, balance.Withdrawal, balance.UserID,
	)
	if err != nil {
		return err
	}

	// проверка успешного обновления
	rows, _ := result.RowsAffected()
	if rows > 0 {
		return nil
	}

	// добавляем запись пользователя
	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO balance (user_id, accrual, withdrawal) VALUES ($1, $2, $3)`,
		balance.UserID, balance.Accrual, balance.Withdrawal,
	)

	return err
}
