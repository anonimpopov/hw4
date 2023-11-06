package authsvc

import (
	"context"
	"fmt"
	"time"

	"github.com/anonimpopov/hw4/internal/model"
	"github.com/anonimpopov/hw4/internal/repo"
	"github.com/anonimpopov/hw4/internal/service"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims

	Login string `json:"login"`
}

type authService struct {
	repo repo.User

	signingKey           string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func (s *authService) makeToken(login string, duration time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
		Login: login,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(s.signingKey))
}

func (s *authService) newTokenPair(login string) (*model.TokenPair, error) {
	accessToken, err := s.makeToken(login, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.makeToken(login, s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &model.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *authService) parseTokenString(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(s.signingKey), nil
	})
}

func (s *authService) Login(ctx context.Context, login, password string) (*model.TokenPair, error) {
	user, err := s.repo.ValidateUser(ctx, login, password)
	if err != nil {
		return nil, err
	}

	return s.newTokenPair(user.Login)
}

func (s *authService) ValidateAndRefresh(ctx context.Context, tokenPair *model.TokenPair) (new *model.TokenPair, err error) {
	accessToken, err := s.parseTokenString(tokenPair.AccessToken)

	switch v := err.(type) {
	case nil:
		return tokenPair, nil

	case *jwt.ValidationError:
		if v.Errors&jwt.ValidationErrorExpired == 0 {
			return nil, err
		}

		_, err = s.parseTokenString(tokenPair.RefreshToken)
		if err != nil {
			return nil, fmt.Errorf("%w: refresh token not valid: %s", service.ErrForbidden, err)
		}

		claims, ok := accessToken.Claims.(*Claims)
		if !ok {
			return nil, service.ErrUnsupportedClaims
		}

		return s.newTokenPair(claims.Login)
	}

	return nil, err
}

func New(config *service.AuthConfig, repo repo.User) service.Auth {
	return &authService{
		repo:                 repo,
		signingKey:           config.SigningKey,
		accessTokenDuration:  config.AccessTokenDuration,
		refreshTokenDuration: config.RefreshTokenDuration,
	}
}
