package profiles

import (
	"context"
	"cvital/db"
	"fmt"
	"log"
)

type useCase struct {
	db db.PostgresDB
}

type UseCase interface {
	CreateCVProfile(ctx context.Context, req CreateCVProfileRequest, cvitalUserID uint) (*CVProfile, error)
}

func NewUseCase(db db.PostgresDB) UseCase {
	return &useCase{
		db: db,
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

func (u *useCase) CreateCVProfile(ctx context.Context, req CreateCVProfileRequest, cvitalUserID uint) (*CVProfile, error) {
	dbRequest := db.CreateCVProfileRequest{
		CvitalUserID:   cvitalUserID,
		CVText:         req.CVText,
		FirstNames:     req.FirstNames,
		Surname:        req.Surname,
		IDNumber:       req.IDNumber,
		PassportNumber: req.PassportNumber,
	}

	newCVProfile, err := u.db.CreateCVProfile(ctx, dbRequest)
	if err != nil {
		return nil, err
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

func (u *useCase) GetUserCVProfile(ctx context.Context, userID uint) (*CVProfile, error) {

	storedProfile, err := u.db.GetCVProfileByUserID(ctx, userID)
	if err != nil {
		log.Printf("GetCVProfileByUserID error: %v", err)
		return nil, fmt.Errorf("Internal Error")
	}

	cvProfile := CVProfile(*storedProfile)

	return &cvProfile, nil
}
