package users

import "fmt"

type User struct {
	FullName          string
	EncryptedPassword string
	EmailAddress      string
}

type CreateUserRequest struct {
	FullName     string
	Password     string
	EmailAddress string
}

func CreateUser(req CreateUserRequest) error {

	return nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(req LoginRequest) error {
	//Stub
	if req.Email == "admin@email.com" && req.Password == "1234abcd" {
		return nil
	}
	return fmt.Errorf("user does not exist")
}
