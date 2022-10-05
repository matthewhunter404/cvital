package profiles

import (
	"context"
	"cvital/db"
	"cvital/domain"

	"github.com/rs/zerolog"
)

type useCase struct {
	db     db.PostgresDB
	logger zerolog.Logger
}

type UseCase interface {
	CreateCVProfile(ctx context.Context, req CreateCVProfileRequest, userEmail string) (*CVProfile, error)
	GetUserCVProfile(ctx context.Context, userEmail string) (*CVProfile, error)
}

func NewUseCase(db db.PostgresDB, logger zerolog.Logger) UseCase {
	return &useCase{
		db:     db,
		logger: logger,
	}
}

type CVProfile struct {
	ID             uint   `json:"id"`
	CvitalUserID   uint   `json:"cvital_user_id"`
	CVText         string `json:"cv_text"`
	FirstNames     string `json:"first_names"`
	Surname        string `json:"surname"`
	IDNumber       string `json:"id_number"`
	PassportNumber string `json:"passport_number"`
}

type CreateCVProfileRequest struct {
	CVText         string `json:"cv_text"`
	FirstNames     string `json:"first_names"`
	Surname        string `json:"surname"`
	IDNumber       string `json:"id_number"`
	PassportNumber string `json:"passport_number"`
}

func (u *useCase) CreateCVProfile(ctx context.Context, req CreateCVProfileRequest, userEmail string) (*CVProfile, error) {

	user, err := u.db.GetUserByEmail(ctx, userEmail)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, domain.ErrNotFound
		default:
			return nil, domain.WrapError(domain.ErrInternal, err)
		}
	}

	dbRequest := db.CreateCVProfileRequest{
		CvitalUserID:   user.ID,
		CVText:         req.CVText,
		FirstNames:     req.FirstNames,
		Surname:        req.Surname,
		IDNumber:       req.IDNumber,
		PassportNumber: req.PassportNumber,
	}

	newCVProfile, err := u.db.CreateCVProfile(ctx, dbRequest)
	if err != nil {
		switch err {
		case db.ErrUniqueViolation:
			return nil, domain.ErrAlreadyExists
		default:
			return nil, domain.WrapError(domain.ErrInternal, err)
		}
	}

	cvProfile := CVProfile{
		ID:             newCVProfile.ID,
		CvitalUserID:   newCVProfile.CvitalUserID,
		CVText:         newCVProfile.CVText,
		FirstNames:     newCVProfile.FirstNames,
		Surname:        newCVProfile.Surname,
		IDNumber:       newCVProfile.IDNumber,
		PassportNumber: newCVProfile.PassportNumber,
	}

	return &cvProfile, nil
}

func (u *useCase) GetUserCVProfile(ctx context.Context, userEmail string) (*CVProfile, error) {

	//TODO, two DB reads?
	user, err := u.db.GetUserByEmail(ctx, userEmail)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, domain.ErrNotFound
		default:
			return nil, domain.WrapError(domain.ErrInternal, err)
		}
	}

	storedProfile, err := u.db.GetCVProfileByUserID(ctx, user.ID)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, domain.ErrNotFound
		default:
			return nil, domain.WrapError(domain.ErrInternal, err)
		}
	}

	cvProfile := CVProfile(*storedProfile)

	return &cvProfile, nil
}
