-- +goose Up
CREATE TABLE cv_user (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(50) NOT NULL,
    encrypted_password VARCHAR(100),
    email_address VARCHAR(50) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE cv_user;