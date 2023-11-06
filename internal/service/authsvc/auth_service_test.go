package authsvc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/anonimpopov/hw4/pkg/jaeger"

	"github.com/anonimpopov/hw4/internal/auth/model"
	"github.com/anonimpopov/hw4/internal/auth/repo/repomock"
	"github.com/anonimpopov/hw4/internal/auth/service"
	"github.com/anonimpopov/hw4/internal/auth/service/authsvc"
)

func TestTokens(t *testing.T) {
	ctx := context.Background()

	userRepo := repomock.NewUser()
	userRepo.On(
		"WithNewTx",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		mock.MatchedBy(func(f func(ctx context.Context) error) bool { return true }),
	).Return(&model.User{Login: "login"}, nil)
	userRepo.On(
		"ValidateUser",
		mock.MatchedBy(func(ctx context.Context) bool { return true }),
		"login",
		"password",
	).Return(&model.User{Login: "login"}, nil)

	config := &service.AuthConfig{
		SigningKey:           "signingKey",
		AccessTokenDuration:  1 * time.Second,
		RefreshTokenDuration: 2 * time.Second,
	}

	svc := authsvc.New(zap.NewNop(), jaeger.NewDummy(), config, userRepo)

	initialPair, err := svc.Login(ctx, "login", "password")
	require.Nil(t, err)

	newPair, err := svc.ValidateAndRefresh(ctx, initialPair)
	require.Nil(t, err)

	require.Equal(t, initialPair.AccessToken, newPair.AccessToken)
	require.Equal(t, initialPair.RefreshToken, newPair.RefreshToken)

	time.Sleep(config.AccessTokenDuration)

	newPair, err = svc.ValidateAndRefresh(ctx, initialPair)
	require.Nil(t, err)

	require.NotEqual(t, initialPair.AccessToken, newPair.AccessToken)
	require.NotEqual(t, initialPair.RefreshToken, newPair.RefreshToken)

	time.Sleep(config.RefreshTokenDuration)

	_, err = svc.ValidateAndRefresh(ctx, newPair)
	require.ErrorIs(t, err, service.ErrForbidden)
}
