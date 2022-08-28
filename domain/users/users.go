package users

import (
	"context"
	"cvital/db"
	"fmt"
	"log"
)

type useCase struct {
	db     db.PostgresDB
	jwtKey string
}

type UseCase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*string, error)
}

func NewUseCase(db db.PostgresDB, jwtKey string) UseCase {
	return &useCase{
		db:     db,
		jwtKey: jwtKey,
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

	user, err := u.db.CreateUser(ctx, dbRequest)
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

	user, err := u.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return nil, fmt.Errorf("login failed")
	}

	passwordCorrect := CheckPasswordHash(req.Password, user.EncryptedPassword)
	if !passwordCorrect {
		return nil, fmt.Errorf("login failed")
	}

	//TODO should this be set using SetCookie on the http response rather than passed back in the body?
	jwt, err := u.CreateJWT(req.Email)
	if err != nil {
		return nil, fmt.Errorf("login failed")
	}
	return jwt, nil
}
