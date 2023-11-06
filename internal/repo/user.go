package repo

import (
	"context"

	"github.com/anonimpopov/hw4/internal/model"
)

type User interface {
	WithNewTx(ctx context.Context, f func(ctx context.Context) error) error
	AddUser(ctx context.Context, login, password, email string) error
	GetUser(ctx context.Context, login string) (*model.User, error)
	ValidateUser(ctx context.Context, login, password string) (*model.User, error)
}
