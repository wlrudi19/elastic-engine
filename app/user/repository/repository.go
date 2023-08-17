package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/wlrudi19/elastic-engine/app/user/model"
)

type UserRepository interface {
	FindUser(ctx context.Context, tx *sql.Tx, email string) (model.UserResponse, error)
}

type userrepository struct {
}

func NewUserRepository() UserRepository {
	return &userrepository{}
}

func (ur *userrepository) FindUser(ctx context.Context, tx *sql.Tx, email string) (model.UserResponse, error) {
	log.Printf("[%s][QUERY] finding user with email: %s", ctx.Value("userEmail"), email)

	var user model.UserResponse

	selectBuilder := squirrel.Select("id, username, created_on").From("users").Where(squirrel.Eq{"email": email}).Where(squirrel.Eq{"deleted_on": nil})
	query, args, err := selectBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		log.Printf("[QUERY] user not found, %v", err)
	}

	err = tx.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Created,
	)

	if err != nil {
		log.Printf("[QUERY] failed to finding user, %v", err)
	}

	return user, err
}
