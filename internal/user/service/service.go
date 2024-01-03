package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/user/model"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserExists = errors.New("логин уже занят")
var ErrInvalidCredentials = errors.New("неверная пара логин/пароль")

type UserRepository interface {
	Create(ctx context.Context, login, password string) (uuid.UUID, error)
	FindByLogin(ctx context.Context, login string) (model.User, error)
}

type UserService struct {
	r UserRepository
}

func NewUserService(repository UserRepository) *UserService {
	return &UserService{r: repository}
}

func (s *UserService) Register(ctx context.Context, login, password string) (uuid.UUID, error) {
	// проверяем наличие пользователя с таким логином
	user, _ := s.r.FindByLogin(ctx, login)
	if user.ID != uuid.Nil {
		return uuid.Nil, ErrUserExists
	}

	// подготавливаем пароль к хранению
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, errors.New("ошибка при создании хэша пароля")
	}

	return s.r.Create(ctx, login, string(passwordHash))
}

func (s *UserService) Login(ctx context.Context, login, password string) (uuid.UUID, error) {
	// проверяем наличие пользователя с таким логином
	user, err := s.r.FindByLogin(ctx, login)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	// проверяем пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	return user.ID, nil
}
