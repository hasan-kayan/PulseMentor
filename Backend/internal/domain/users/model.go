package users

import "time"

type User struct {
	ID        string
	Email     string
	Password  string // hashed
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserInput struct {
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

