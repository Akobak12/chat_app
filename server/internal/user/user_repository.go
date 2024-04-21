package user

import (
	"context"
	"database/sql"
)

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type RepositoryDb struct {
	db DBTX
}

func NewRepository(db DBTX) Repository {
	return &RepositoryDb{db: db}
}

func (repository *RepositoryDb) CreateUser(ctx context.Context, user *User) (*User, error) {
	var createdId uint64

	query := "INSERT INTO public.users(username, password, email) VALUES ($1, $2, $3) returning id"
	err := repository.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Email).Scan(&createdId)
	if err != nil {
		return nil, err
	}

	user.Id = createdId
	return user, nil
}

func (repository *RepositoryDb) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	user := User{}
	query := "SELECT id, email, username, password FROM public.users WHERE email = $1"
	err := repository.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Email, &user.Username, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
