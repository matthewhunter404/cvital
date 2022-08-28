package users

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (u *useCase) CreateJWT(email string) (*string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(u.jwtKey)
	if err != nil {
		return nil, err
	}
	return &tokenString, nil
}
