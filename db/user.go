package db

import (
	"context"
	"log"
)

type User struct {
	ID                uint   `db:"id"`
	FullName          string `db:"full_name"`
	EncryptedPassword string `db:"encrypted_password"`
	EmailAddress      string `db:"email_address"`
}

type CreateUserRequest struct {
	FullName          string
	EncryptedPassword string
	EmailAddress      string
}

func (d PostgresDB) CreateUser(ctx context.Context, req CreateUserRequest) (*User, error) {
	sqlStatement := `INSERT INTO cvital_user (full_name, encrypted_password, email_address) VALUES ($1, $2, $3) RETURNING id`

	var id uint
	err := d.QueryRowContext(ctx, sqlStatement, req.FullName, req.EncryptedPassword, req.EmailAddress).Scan(&id)
	if err != nil {
		return nil, err
	}

	user := User{
		ID:                id,
		FullName:          req.FullName,
		EncryptedPassword: req.EncryptedPassword,
		EmailAddress:      req.EmailAddress,
	}
	return &user, nil
}

func (d PostgresDB) GetUserByEmail(ctx context.Context, emailAddress string) (*User, error) {
	sqlStatement := `SELECT id, full_name, encrypted_password, email_address FROM cvital_user WHERE email_address = $1`
	var user User
	err := d.GetContext(ctx, &user, sqlStatement, emailAddress)
	if err != nil {
		return nil, err
	}
	log.Printf("user: %v", user)
	return &user, nil
}
