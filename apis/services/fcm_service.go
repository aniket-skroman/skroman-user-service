package services

import (
	"github.com/aniket-skroman/skroman-user-service/apis/dtos"
	"github.com/aniket-skroman/skroman-user-service/apis/helper"
	"github.com/aniket-skroman/skroman-user-service/apis/repositories"
	db "github.com/aniket-skroman/skroman-user-service/sqlc_lib"
	"github.com/google/uuid"
)

type FCMService interface {
	CreateUserFCMData(req dtos.CreateUserFCMDataRequestDTO) (db.UserFcmData, error)
	FetchTokenByUsers(user_id string) ([]db.UserFcmData, error)
}

type fcm_serv struct {
	fcm_repo repositories.FCMRepository
}

func NewFCMService(repo repositories.FCMRepository) FCMService {
	return &fcm_serv{
		fcm_repo: repo,
	}
}

func (serv *fcm_serv) CreateUserFCMData(req dtos.CreateUserFCMDataRequestDTO) (db.UserFcmData, error) {
	user_obj_id, err := uuid.Parse(req.UserID)

	if err != nil {
		return db.UserFcmData{}, helper.ERR_INVALID_ID
	}

	args := db.CreateUserFCMDataParams{
		UserID:   user_obj_id,
		FcmToken: req.FcmToken,
	}

	result, err := serv.fcm_repo.CreateUserFCMData(args)

	err = helper.Handle_DBError(err)

	return result, err
}

func (serv *fcm_serv) FetchTokenByUsers(user_id string) ([]db.UserFcmData, error) {
	user_obj_id, err := uuid.Parse(user_id)

	if err != nil {
		return nil, helper.ERR_INVALID_ID
	}

	result, err := serv.fcm_repo.FetchTokenByUsers(user_obj_id)

	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, helper.Err_Data_Not_Found
	}

	return result, nil
}
