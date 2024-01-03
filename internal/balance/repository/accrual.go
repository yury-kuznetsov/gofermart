package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/balance/model"
	"time"
)

type AccrualRepository struct {
	db *sql.DB
}

func NewAccrualRepository(db *sql.DB) *AccrualRepository {
	r := &AccrualRepository{db: db}

	_, _ = r.db.Exec(`CREATE TABLE IF NOT EXISTS balance_accrual (
		id         uuid    not null constraint orders_pk primary key,
		user_id    uuid    not null constraint orders_users_id_fk references "user",
		number     varchar not null constraint orders_pk_2 unique,
		status     varchar not null,
		sum    	   decimal,
		created_at timestamp
	)`)

	return r
}

func (r *AccrualRepository) Save(ctx context.Context, model model.Accrual) error {
	queryUpdate := "UPDATE balance_accrual SET status = $1, sum = $2 WHERE id = $3"
	queryInsert := "INSERT INTO balance_accrual VALUES ($1, $2, $3, $4, $5, $6)"

	// обновляем запись
	result, err := r.db.ExecContext(
		ctx,
		queryUpdate,
		model.Status, model.Sum, model.ID,
	)
	if err != nil {
		return err
	}

	// проверка успешного обновления
	rows, _ := result.RowsAffected()
	if rows > 0 {
		return nil
	}

	// добавляем запись
	_, err = r.db.ExecContext(
		ctx,
		queryInsert,
		model.ID, model.UserID, model.Number, model.Status, model.Sum, model.CreatedAt.Format(time.RFC3339),
	)

	return err
}

func (r *AccrualRepository) FindByNumber(ctx context.Context, number string) (model.Accrual, error) {
	var accrual model.Accrual
	err := r.db.QueryRowContext(
		ctx,
		"SELECT * FROM balance_accrual WHERE number = $1",
		number,
	).Scan(&accrual.ID, &accrual.UserID, &accrual.Number, &accrual.Status, &accrual.Sum, &accrual.CreatedAt)

	return accrual, err
}

func (r *AccrualRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]model.Accrual, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT * FROM balance_accrual WHERE user_id = $1 ORDER BY created_at",
		userID,
	)
	if err != nil {
		return nil, err
	}

	return parseRows(rows)
}

func (r *AccrualRepository) FindForSync(ctx context.Context) ([]model.Accrual, error) {
	rows, err := r.db.QueryContext(
		ctx,
		"SELECT * FROM balance_accrual WHERE status = $1",
		model.StatusNew,
	)
	if err != nil {
		return nil, err
	}

	return parseRows(rows)
}

func parseRows(rows *sql.Rows) ([]model.Accrual, error) {
	var accruals []model.Accrual
	for rows.Next() {
		var accrual model.Accrual
		err := rows.Scan(
			&accrual.ID,
			&accrual.UserID,
			&accrual.Number,
			&accrual.Status,
			&accrual.Sum,
			&accrual.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		accruals = append(accruals, accrual)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accruals, nil
}
