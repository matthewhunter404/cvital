package users

import (
	"context"
	"cvital/db"
	"cvital/domain"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

type useCase struct {
	db     db.CVitalDB
	jwtKey string
	logger zerolog.Logger
}

type UseCase interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)
	Login(ctx context.Context, req LoginRequest) (*string, *time.Time, error)
	ValidateToken(tokenString string) (*Claims, error)
	CreateJWT(email string) (*string, *time.Time, error)
}

func NewUseCase(db db.CVitalDB, jwtKey string, logger zerolog.Logger) UseCase {
	return &useCase{
		db:     db,
		jwtKey: jwtKey,
		logger: logger,
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
		return nil, domain.WrapError(domain.ErrInternal, err)
	}

	dbRequest := db.CreateUserRequest{
		FullName:          req.FullName,
		EncryptedPassword: hashedPassword,
		EmailAddress:      req.Email,
	}

	user, err := u.db.CreateUser(ctx, dbRequest)
	if err != nil {
		switch err {
		case db.ErrUniqueViolation:
			return nil, domain.ErrAlreadyExists
		default:
			return nil, domain.WrapError(domain.ErrInternal, err)
		}
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
		u.logger.Printf("GetUserByEmail error: %v", err)
		return nil, nil, domain.WrapError(domain.ErrLoginFailed, err)
	}

	passwordCorrect := CheckPasswordHash(req.Password, user.EncryptedPassword)
	if !passwordCorrect {
		u.logger.Printf("Password incorrect\n")
		return nil, nil, domain.WrapError(domain.ErrLoginFailed, fmt.Errorf("Invalid Password"))
	}

	jwt, expiryTime, err := u.CreateJWT(req.Email)
	if err != nil {
		u.logger.Printf("CreateJWT failed: %v \n", err)
		return nil, nil, domain.WrapError(domain.ErrLoginFailed, err)
	}
	return jwt, expiryTime, nil
}
