package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/wlrudi19/elastic-engine/app/user/model"
	"github.com/wlrudi19/elastic-engine/app/user/repository"
)

type UserLogic interface {
	FindUserLogic(ctx context.Context, email string) (model.UserResponse, error)
}

type userlogic struct {
	UserRepository repository.UserRepository
	db             *sql.DB
}

func NewUserLogic(userRepository repository.UserRepository, db *sql.DB) UserLogic {
	return &userlogic{
		UserRepository: userRepository,
		db:             db,
	}
}

func (l *userlogic) FindUserLogic(ctx context.Context, email string) (model.UserResponse, error) {
	log.Printf("[%s][LOGIC] find user with email: %s", ctx.Value("userEmail"), email)

	var user model.UserResponse

	tx, err := l.db.Begin()

	if err != nil {
		log.Printf("[LOGIC] failed to find user :%v", err)
		return user, err
	}

	user, err = l.UserRepository.FindUser(ctx, tx, email)

	if err != nil {
		log.Printf("[LOGIC] failed to find user :%v", err)
		tx.Rollback()
		return user, err
	}

	tx.Commit()
	log.Printf("[%s][LOGIC] user find successfulyy, email: %s", ctx.Value("userEmail"), email)
	return user, nil
}
