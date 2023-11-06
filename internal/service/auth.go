package service

import (
	"context"
	"time"

	"github.com/anonimpopov/hw4/internal/model"
)

type AuthConfigUser struct {
	Login   string `yaml:"login"`
	Pasword string `yaml:"password"`
	Email   string `yaml:"email"`
}

type AuthConfig struct {
	SigningKey           string           `yaml:"signing_key"`
	AccessTokenDuration  time.Duration    `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration    `yaml:"refresh_token_duration"`
	Users                []AuthConfigUser `yaml:"users"`
}

type Auth interface {
	Login(ctx context.Context, login, password string) (*model.TokenPair, error)
	ValidateAndRefresh(ctx context.Context, tokenPair *model.TokenPair) (new *model.TokenPair, err error)
}
