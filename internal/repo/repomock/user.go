package repomock

import (
	"context"

	"github.com/anonimpopov/hw4/internal/model"
	"github.com/anonimpopov/hw4/internal/repo"

	"github.com/stretchr/testify/mock"
)

var _ repo.User = &UserMock{}

type UserMock struct {
	mock.Mock
}

func (m *UserMock) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	args := m.Called(ctx, f)
	return args.Error(0)
}

func (m *UserMock) AddUser(ctx context.Context, login, password, email string) error {
	args := m.Called(ctx, login, password, email)
	return args.Error(0)
}

func (m *UserMock) GetUser(ctx context.Context, login string) (*model.User, error) {
	args := m.Called(ctx, login)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *UserMock) ValidateUser(ctx context.Context, login, password string) (*model.User, error) {
	args := m.Called(ctx, login, password)
	return args.Get(0).(*model.User), args.Error(1)
}

func NewUser() *UserMock {
	return &UserMock{}
}
