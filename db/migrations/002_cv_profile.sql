-- +goose Up
CREATE TABLE cv_profile (
    id SERIAL PRIMARY KEY,
    cvital_user_id INT,
    cv_text TEXT,
    first_names VARCHAR(50),
    surname VARCHAR(50),
    id_number VARCHAR(50),
    passport_number VARCHAR(9),
    CONSTRAINT fk_cvital_user FOREIGN KEY(cvital_user_id) REFERENCES cvital_user(id)
);

-- +goose Down
DROP TABLE cv_profile;