package service

import (
	"context"
	"database/sql"
	"log"

	"github.com/wlrudi19/elastic-engine/app/user/model"
	"github.com/wlrudi19/elastic-engine/app/user/repository"
	"github.com/wlrudi19/elastic-engine/helper/jwt"
)

type UserLogic interface {
	FindUserLogic(ctx context.Context, email string) (model.UserResponse, error)
	LoginUserLogic(ctx context.Context, email string) (model.LoginResponse, error)
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

func (l *userlogic) LoginUserLogic(ctx context.Context, email string) (model.LoginResponse, error) {
	log.Printf("[%s][LOGIC] login with email: %s", ctx.Value("loginEmail"), email)

	var login model.LoginResponse

	user, err := l.FindUserLogic(ctx, email)

	if err != nil {
		log.Printf("[LOGIC] failed to find user, %v", err)
		return login, err
	}

	token, err := jwt.NewJWT().GenerateAccessToken(user.Id, email)

	if err != nil {
		log.Printf("[LOGIC] failed to generate access token, %v", err)
		return login, err
	}

	login = model.LoginResponse{
		Id:          user.Id,
		AccessToken: token,
	}

	log.Printf("[%s][LOGIC] login successfulyy, with token: %s", ctx.Value("loginToken"), token)
	return login, nil
}
