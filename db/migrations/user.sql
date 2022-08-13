-- +goose Up
CREATE TABLE user (
    id SERIAL PRIMARY KEY,
    fullname VARCHAR(50) NOT NULL,
    encrypted_password VARCHAR(100),
    email_address VARCHAR(50) UNIQUE NOT NULL
);


-- +goose Down
DROP TABLE user;