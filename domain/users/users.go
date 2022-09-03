package users

import (
	"context"
	"cvital/db"
	"fmt"
	"log"
	"time"
)

type useCase struct {
	db     db.PostgresDB
	jwtKey string
}

type UseCase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*string, *time.Time, error)
	ValidateToken(tokenString string) (*Claims, error)
	CreateJWT(email string) (*string, *time.Time, error)
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

func (u *useCase) Login(ctx context.Context, req LoginRequest) (*string, *time.Time, error) {

	user, err := u.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return nil, nil, fmt.Errorf("login failed")
	}

	passwordCorrect := CheckPasswordHash(req.Password, user.EncryptedPassword)
	if !passwordCorrect {
		log.Printf("Password incorrect: %v \n", err)
		return nil, nil, fmt.Errorf("login failed")
	}

	jwt, expiyTime, err := u.CreateJWT(req.Email)
	if err != nil {
		log.Printf("CreateJWT failed: %v \n", err)
		return nil, nil, fmt.Errorf("login failed")
	}
	return jwt, expiyTime, nil
}
