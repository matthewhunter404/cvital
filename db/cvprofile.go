package db

import (
	"context"
)

type CVProfile struct {
	ID             uint   `db:"id"`
	CvitalUserID   uint   `db:"cvital_user_id"`
	CVText         string `db:"cv_text"`
	FirstNames     string `db:"first_names"`
	Surname        string `db:"surname"`
	IDNumber       string `db:"id_number"`
	PassportNumber string `db:"passport_number"`
}

type CreateCVProfile struct {
	CvitalUserID   uint
	CVText         string
	FirstNames     string
	Surname        string
	IDNumber       string
	PassportNumber string
}

func (d PostgresDB) CreateCVProfile(ctx context.Context, req CreateCVProfile) (*CVProfile, error) {
	sqlStatement := `INSERT INTO cv_profile (cvital_user_id, cv_text, first_names, surname, id_number, passport_number) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id uint
	err := d.QueryRowContext(ctx, sqlStatement, req.CvitalUserID, req.CVText, req.FirstNames, req.Surname, req.IDNumber, req.PassportNumber).Scan(&id)
	if err != nil {
		return nil, err
	}

	cvProfile := CVProfile{
		ID:             id,
		CVText:         req.CVText,
		FirstNames:     req.FirstNames,
		Surname:        req.Surname,
		IDNumber:       req.IDNumber,
		PassportNumber: req.PassportNumber,
	}
	return &cvProfile, nil
}

func (d PostgresDB) GetCVProfileByUserID(ctx context.Context, cvitalUserID uint) (*CVProfile, error) {
	sqlStatement := `SELECT id, cvital_user_id, cv_text, first_names, surname, id_number, passport_number FROM cv_profile WHERE cvital_user_id = $1`
	var cvProfile CVProfile
	err := d.GetContext(ctx, &cvProfile, sqlStatement, cvitalUserID)
	if err != nil {
		return nil, err
	}
	return &cvProfile, nil
}
