package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/user/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	r := &UserRepository{db: db}

	_, _ = r.db.Exec(`CREATE TABLE IF NOT EXISTS "user" (
		id       uuid    not null constraint users_pk primary key,
		login    varchar not null constraint users_pk_2 unique,
		password varchar not null
	)`)

	return r
}

func (r *UserRepository) Create(ctx context.Context, login, password string) (uuid.UUID, error) {
	id := uuid.New()
	query := `INSERT INTO "user" (id, login, password) VALUES ($1, $2, $3)`
	if _, err := r.db.ExecContext(ctx, query, id, login, password); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (r *UserRepository) FindByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, login, password FROM "user" WHERE login = $1`,
		login,
	).Scan(&user.ID, &user.Login, &user.Password)

	return user, err
}
