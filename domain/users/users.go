package users

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
	FullName     string
	Password     string
	EmailAddress string
}

func Login(req LoginRequest) error {

	return nil
}
