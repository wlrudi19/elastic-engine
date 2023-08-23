package service

import (
	"context"
	"log"

	"github.com/wlrudi19/elastic-engine/app/user/model"
	"github.com/wlrudi19/elastic-engine/app/user/repository"
	"github.com/wlrudi19/elastic-engine/infrastructure/middlewares"
)

type UserLogic interface {
	FindUserLogic(ctx context.Context, email string) (model.UserResponse, error)
	LoginUserLogic(ctx context.Context, email string) (model.LoginResponse, error)
}

type userlogic struct {
	UserRepository repository.UserRepository
}

func NewUserLogic(userRepository repository.UserRepository) UserLogic {
	return &userlogic{
		UserRepository: userRepository,
	}
}

func (l *userlogic) FindUserLogic(ctx context.Context, email string) (model.UserResponse, error) {
	log.Printf("[LOGIC] find user with email: %s", email)

	var user model.UserResponse

	user, err := l.UserRepository.FindUser(ctx, email)

	if err != nil {
		log.Printf("[LOGIC] failed to find user :%v", err)
		return user, err
	}

	log.Printf("[LOGIC] user find successfulyy, email: %s", email)
	return user, nil
}

func (l *userlogic) LoginUserLogic(ctx context.Context, email string) (model.LoginResponse, error) {
	log.Printf("[LOGIC] login with email: %s", email)

	var login model.LoginResponse

	user, err := l.FindUserLogic(ctx, email)

	if err != nil {
		log.Printf("[LOGIC] failed to find user, %v", err)
		return login, err
	}

	token, err := middlewares.GenerateAccessToken(user.Id, email)

	if err != nil {
		log.Printf("[LOGIC] failed to generate access token, %v", err)
		return login, err
	}

	login = model.LoginResponse{
		Id:          user.Id,
		AccessToken: token,
	}

	log.Printf("[LOGIC] login successfulyy, with token: %s", token)
	return login, nil
}
