package repository

import (
	"context"
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"github.com/wlrudi19/elastic-engine/app/user/model"
)

type UserRepository interface {
	FindUser(ctx context.Context, email string) (model.UserResponse, error)
	FindUserRedis(ctx context.Context, email string) (model.UserResponse, error)
	WithTransaction() (*sql.Tx, error)
	GetUserRedis(ctx context.Context, email string) (model.UserResponse, error)
	SetUserRedis(ctx context.Context, email string, user model.UserResponse) error
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

func (ur *userrepository) GetUserRedis(ctx context.Context, email string) (model.UserResponse, error) {
	var userCache model.UserResponse

	userCacheKey := "user:" + email
	userFields, err := ur.redis.HGetAll(ctx, userCacheKey).Result()

	if err != nil {
		return userCache, err
	}

	if len(userFields) == 0 {
		return userCache, redis.Nil
	}

	id, _ := strconv.Atoi(userFields["id"])
	name := userFields["name"]
	created, _ := time.Parse(time.RFC3339, userFields["created"])

	userCache = model.UserResponse{
		Id:      id,
		Name:    name,
		Created: created,
	}

	return userCache, nil
}

func (ur *userrepository) SetUserRedis(ctx context.Context, email string, user model.UserResponse) error {
	userCacheKey := "user:" + email

	pipe := ur.redis.TxPipeline()
	defer pipe.Close()

	pipe.HMSet(ctx, userCacheKey, map[string]interface{}{
		"id":      user.Id,
		"name":    user.Name,
		"created": user.Created,
	})

	pipe.Expire(ctx, userCacheKey, 1*time.Hour)
	_, err := pipe.Exec(ctx)

	if err != nil {
		pipe.Discard()
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

	//userCacheJSON, _ := json.Marshal(user)
	err = ur.SetUserRedis(ctx, email, user)

	if err != nil {
		log.Printf("[QUERY] failed to cache user data: %v", err)
		return user, err
	}

	return user, err
}

func (ur *userrepository) FindUserRedis(ctx context.Context, email string) (model.UserResponse, error) {
	log.Printf("[REDIS] finding user redis with email: %s", email)

	//var userCache model.UserResponse

	userCache, err := ur.GetUserRedis(ctx, email)

	if err != nil {
		log.Printf("[REDIS] user not found in redis, %v", err)
		return userCache, err
	}

	// err = json.Unmarshal(userCacheJSON, &userCache)

	// if err != nil {
	// 	log.Printf("[REDIS] error unmarshalling user: %v", err)
	// 	return userCache, err
	// }

	log.Printf("[REDIS] user data found in redis, email: %s", email)
	return userCache, nil
}
