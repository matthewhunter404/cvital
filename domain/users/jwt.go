package users

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Create the JWT key used to create the signature

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (u *useCase) CreateJWT(email string) (*string, *time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(([]byte(u.jwtKey)))
	if err != nil {
		return nil, nil, err
	}

	return &tokenString, &expirationTime, nil
}

func (u *useCase) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("jwt error: token not valid")
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		return nil, fmt.Errorf("jwt error: token expired")
	}

	return claims, nil
}

func GetTokenFromBearerAuth(bearerAuth string) (string, error) {
	if !(bearerAuth[:7] == "Bearer ") {
		return "", fmt.Errorf("invalid bearer token format")
	}
	return bearerAuth[7:], nil
}
