-- +goose Up
CREATE TABLE cv_profile (
    id SERIAL PRIMARY KEY,
    cvital_user_id INT,
    cv_text TEXT,
    first_names VARCHAR(50),
    surname VARCHAR(50),
    id_number VARCHAR(50),
    passport_number VARCHAR(9),
    CONSTRAINT cv_profile_cvital_user_id_fkey FOREIGN KEY(cvital_user_id) REFERENCES cvital_user(id) ON DELETE CASCADE,
    CONSTRAINT cv_profile_cvital_user_id_key UNIQUE(cvital_user_id)
);

-- +goose Down
DROP TABLE cv_profile;