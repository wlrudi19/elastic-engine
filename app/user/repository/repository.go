package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/wlrudi19/elastic-engine/app/user/model"
)

type UserRepository interface {
	FindUser(ctx context.Context, email string) (model.UserResponse, error)
	FindUserRedis(ctx context.Context, email string) (model.UserResponse, error)
	WithTransaction() (*sql.Tx, error)
	GetUserRedis(ctx context.Context, email string) ([]byte, error)
	SetUserRedis(ctx context.Context, email string, json []byte) error
}

type userrepository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewUserRepository(db *sql.DB, redis *redis.Client) UserRepository {
	return &userrepository{
		db:    db,
		redis: redis,
	}
}

func (pr *userrepository) WithTransaction() (*sql.Tx, error) {
	tx, err := pr.db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (ur *userrepository) GetUserRedis(ctx context.Context, email string) ([]byte, error) {
	userCacheKey := "user:" + email
	userCacheJSON, err := ur.redis.Get(ctx, userCacheKey).Bytes()
	if err != nil {
		return userCacheJSON, err
	}

	return userCacheJSON, nil
}

func (ur *userrepository) SetUserRedis(ctx context.Context, email string, json []byte) error {
	userCacheKey := "user" + email
	err := ur.redis.Set(ctx, userCacheKey, json, 24*time.Hour).Err()
	if err != nil {
		return err
	}

	return nil
}

func (ur *userrepository) FindUser(ctx context.Context, email string) (model.UserResponse, error) {
	log.Printf("[QUERY] finding user with email: %s", email)

	var user model.UserResponse

	selectBuilder := squirrel.Select("id, username, created_on").From("users").Where(squirrel.Eq{"email": email}).Where(squirrel.Eq{"deleted_on": nil})
	query, args, err := selectBuilder.PlaceholderFormat(squirrel.Dollar).ToSql()

	if err != nil {
		log.Printf("[QUERY] user not found, %v", err)
		return user, err
	}

	err = ur.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id,
		&user.Name,
		&user.Created,
	)

	if err != nil {
		log.Printf("[QUERY] failed to finding user, %v", err)
		return user, err
	}

	userCacheJSON, _ := json.Marshal(user)
	err = ur.SetUserRedis(ctx, email, userCacheJSON)

	if err != nil {
		log.Printf("[QUERY] failed to cache user data: %v", err)
		return user, err
	}

	return user, err
}

func (ur *userrepository) FindUserRedis(ctx context.Context, email string) (model.UserResponse, error) {
	log.Printf("[REDIS] finding user redis with email: %s", email)

	var userCache model.UserResponse

	userCacheJSON, err := ur.GetUserRedis(ctx, email)

	if err != nil {
		log.Printf("[REDIS] user not found in redis, %v", err)
		return userCache, err
	}

	err = json.Unmarshal(userCacheJSON, &userCache)

	if err != nil {
		log.Printf("[REDIS] error unmarshalling user: %v", err)
		return userCache, err
	}

	log.Printf("[REDIS] user data found in redis, email: %s", email)
	return userCache, nil
}
