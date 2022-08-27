package users

import (
	"context"
	"cvital/db"
	"fmt"
)

type useCase struct {
	DB db.PostgresDB
}

type UseCase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) error
}

func NewUseCase(db db.PostgresDB) UseCase {
	return &useCase{
		DB: db,
	}
}

type User struct {
	FullName          string `json:"full_name"`
	EncryptedPassword string `json:"-"`
	EmailAddress      string `json:"email"`
}

type CreateUserRequest struct {
	FullName string `json:"full_name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (u *useCase) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {

	dbRequest := db.CreateUserRequest{
		FullName:          req.FullName,
		EncryptedPassword: req.Password,
		EmailAddress:      req.Email,
	}
	user, err := u.DB.CreateUser(ctx, dbRequest)
	if err != nil {
		return nil, err
	}

	newUser := User{
		FullName:          user.FullName,
		EncryptedPassword: user.EncryptedPassword,
		EmailAddress:      user.EmailAddress,
	}

	return &newUser, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *useCase) Login(ctx context.Context, req LoginRequest) error {
	//Stub
	if req.Email == "admin@email.com" && req.Password == "1234abcd" {
		return nil
	}
	return fmt.Errorf("user does not exist")
}
