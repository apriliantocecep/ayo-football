package model

import "time"

type LoginRequest struct {
	Identity string
	Password string
}

type LoginResponse struct {
	AccessToken          string
	AccessTokenExpiresAt time.Time
}

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterResponse struct {
	UserId   string
	Username string
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Username string
	Password string
}

type CreateUserResponse struct {
	UserId string
}
