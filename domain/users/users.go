package users

import (
	"context"
	"cvital/db"
	"fmt"
	"log"
)

type useCase struct {
	DB db.PostgresDB
}

type UseCase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*string, error)
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

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	dbRequest := db.CreateUserRequest{
		FullName:          req.FullName,
		EncryptedPassword: hashedPassword,
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

func (u *useCase) Login(ctx context.Context, req LoginRequest) (*string, error) {

	user, err := u.DB.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("User login request failed as user does not exist", err)
		return nil, fmt.Errorf("login failed")
	}

	passwordCorrect := CheckPasswordHash(req.Password, user.EncryptedPassword)
	if !passwordCorrect {
		return nil, fmt.Errorf("login failed")
	}
	JWTStub := "LoginSuccessfulJWT"
	return &JWTStub, nil
}
