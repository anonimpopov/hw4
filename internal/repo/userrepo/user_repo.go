package userrepo

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/anonimpopov/hw4/internal/model"
	"github.com/anonimpopov/hw4/internal/repo"
	"github.com/anonimpopov/hw4/internal/service"
)

type userRepo struct {
	pgxPool *pgxpool.Pool
}

func (r *userRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *userRepo) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	return r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(context.WithValue(ctx, repo.CtxKeyTx, tx))
	})
}

func (r *userRepo) AddUser(ctx context.Context, login, password, email string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = r.conn(ctx).Exec(ctx, `INSERT INTO users (login, password_hash, email) VALUES ($1, $2, $3)`, login, hash, email)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetUser(ctx context.Context, login string) (*model.User, error) {
	var user model.User

	row := r.conn(ctx).QueryRow(ctx, `SELECT login, password_hash, email FROM users WHERE login = $1`, login)
	if err := row.Scan(&user.Login, &user.HashedPassword, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return nil, service.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *userRepo) ValidateUser(ctx context.Context, login, password string) (*model.User, error) {
	user, err := r.GetUser(ctx, login)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		return nil, service.ErrForbidden
	}

	return user, nil
}

func New(config *service.AuthConfig, pgxPool *pgxpool.Pool) (repo.User, error) {
	r := &userRepo{
		pgxPool: pgxPool,
	}

	ctx := context.Background()

	err := r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		for _, user := range config.Users {
			if err := r.AddUser(ctx, user.Login, user.Pasword, user.Email); err != nil {
				log.Fatal(err.Error())
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}
