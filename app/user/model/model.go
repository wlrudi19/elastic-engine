package model

import "time"

type UserResponse struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

type UserRequest struct {
	Email string `json:"email"`
}

type LoginResponse struct {
	AccessToken string `json:"access-token"`
}
