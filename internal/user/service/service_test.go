package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/yury-kuznetsov/gofermart/internal/user/model"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

type mockUserRepository struct {
	users []model.User
}

func (m *mockUserRepository) FindByLogin(_ context.Context, login string) (model.User, error) {
	for _, user := range m.users {
		if user.Login == login {
			return user, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func (m *mockUserRepository) Create(_ context.Context, login, password string) (uuid.UUID, error) {
	id := uuid.New()
	user := model.User{ID: id, Login: login, Password: password}
	m.users = append(m.users, user)
	return id, nil
}

func createPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func TestRegister(t *testing.T) {
	repo := &mockUserRepository{
		users: []model.User{},
	}
	svc := UserService{r: repo}

	testCases := []struct {
		name          string
		login         string
		password      string
		expectedError error
	}{
		{name: "SuccessfulRegistration", login: "test", password: "password", expectedError: nil},
		{name: "DuplicateRegistration", login: "test", password: "password", expectedError: ErrUserExists},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), testCase.login, testCase.password)
			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("expected error: %v, but got: %v", testCase.expectedError, err)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	repo := &mockUserRepository{
		users: []model.User{
			{ID: uuid.New(), Login: "admin", Password: createPassword("123")},
		},
	}
	svc := UserService{r: repo}

	testCases := []struct {
		name          string
		login         string
		password      string
		expectedError error
	}{
		{name: "CorrectLogin", login: "admin", password: "123", expectedError: nil},
		{name: "WrongPassword", login: "admin", password: "wrong", expectedError: ErrInvalidCredentials},
		{name: "UserNotExists", login: "guest", password: "123", expectedError: ErrInvalidCredentials},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := svc.Login(context.Background(), testCase.login, testCase.password)
			if !errors.Is(err, testCase.expectedError) {
				t.Errorf("expected error: %v, but got: %v", testCase.expectedError, err)
			}
		})
	}
}
