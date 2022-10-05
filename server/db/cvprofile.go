package db

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
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

// TODO, combine requests?
type CreateCVProfileRequest struct {
	CvitalUserID   uint
	CVText         string
	FirstNames     string
	Surname        string
	IDNumber       string
	PassportNumber string
}

type UpdateCVProfileRequest struct {
	CvitalUserID   uint
	CVText         string
	FirstNames     string
	Surname        string
	IDNumber       string
	PassportNumber string
}

func (d *PostgresDB) CreateCVProfile(ctx context.Context, req CreateCVProfileRequest) (*CVProfile, error) {
	sqlStatement := `INSERT INTO cv_profile (cvital_user_id, cv_text, first_names, surname, id_number, passport_number) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`

	var id uint
	err := d.sqlxDB.QueryRowContext(ctx, sqlStatement, req.CvitalUserID, req.CVText, req.FirstNames, req.Surname, req.IDNumber, req.PassportNumber).Scan(&id)
	if err != nil {
		if sqlxDBErr, ok := err.(*pq.Error); ok {
			if sqlxDBErr.Code.Name() == "unique_violation" {
				return nil, ErrUniqueViolation
			}
		}
		return nil, WrapError(ErrInternal, err)
	}

	cvProfile := CVProfile{
		ID:             id,
		CvitalUserID:   req.CvitalUserID,
		CVText:         req.CVText,
		FirstNames:     req.FirstNames,
		Surname:        req.Surname,
		IDNumber:       req.IDNumber,
		PassportNumber: req.PassportNumber,
	}
	return &cvProfile, nil
}

func (d *PostgresDB) GetCVProfileByUserID(ctx context.Context, cvitalUserID uint) (*CVProfile, error) {
	sqlStatement := `SELECT id, cvital_user_id, cv_text, first_names, surname, id_number, passport_number FROM cv_profile WHERE cvital_user_id = $1`
	var cvProfile CVProfile
	err := d.sqlxDB.GetContext(ctx, &cvProfile, sqlStatement, cvitalUserID)
	if err != nil {
		if sqlxDBErr, ok := err.(*pq.Error); ok {
			if sqlxDBErr.Code.Name() == "case_not_found" {
				return nil, ErrNotFound
			}
		}
		return nil, WrapError(ErrInternal, err)
	}
	return &cvProfile, nil
}

func (d *PostgresDB) UpdateCVProfile(ctx context.Context, req UpdateCVProfileRequest) (*CVProfile, error) {
	d.logger.Debug().Interface("UpdateCVProfileRequest", req).Msg("")
	sqlStatement := `UPDATE cv_profile SET cv_text = $1, first_names = $2, surname = $3, id_number = $4, passport_number = $5 WHERE cvital_user_id = $6 RETURNING id`

	var id uint
	err := d.sqlxDB.QueryRowContext(ctx, sqlStatement, req.CVText, req.FirstNames, req.Surname, req.IDNumber, req.PassportNumber, req.CvitalUserID).Scan(&id)
	if err != nil {
		if sqlxDBErr, ok := err.(*pq.Error); ok {
			if sqlxDBErr.Code.Name() == "unique_violation" {
				return nil, ErrUniqueViolation
			}
		}
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, WrapError(ErrInternal, err)
	}

	cvProfile := CVProfile{
		ID:             id,
		CvitalUserID:   req.CvitalUserID,
		CVText:         req.CVText,
		FirstNames:     req.FirstNames,
		Surname:        req.Surname,
		IDNumber:       req.IDNumber,
		PassportNumber: req.PassportNumber,
	}
	return &cvProfile, nil
}
