package service

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	app_model "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/model"
	app_option "github.com/mbilarusdev/durak_auth_bot/internal/structs/app/option"
	"github.com/mbilarusdev/jwt/jwt"
	jwtmodels "github.com/mbilarusdev/jwt/models"
)

type TokenManager interface {
	FindActualByToken(token string) (*app_model.Token, error)
	FindActualByPlayerID(playerID uint64) (*app_model.Token, error)
	IssueToken(playerID uint64) (*app_model.Token, error)
	BlockToken(tokenID uint64) error
}

type TokenService struct {
	tokenRepository repository.TokenProvider
}

func NewTokenService(tokenRepository repository.TokenProvider) *TokenService {
	service := new(TokenService)
	service.tokenRepository = tokenRepository
	return service
}

func (service *TokenService) FindActualByToken(token string) (*app_model.Token, error) {
	playerID := jwt.GetSubID(token, common.Conf.Token)
	return service.FindActualByPlayerID(playerID)
}

func (service *TokenService) FindActualByPlayerID(playerID uint64) (*app_model.Token, error) {
	finded, err := service.tokenRepository.FindOne(&app_option.TokenFindOptions{PlayerID: playerID})
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	available := jwt.Check(finded.Jwt, common.Conf.SecretKey)
	if !available && finded.Status != app_model.TokenBlocked {
		err := service.tokenRepository.UpdateStatus(finded.ID, app_model.TokenExpired)
		if err != nil {
			return nil, err
		}
		token, err := service.tokenRepository.FindOne(
			&app_option.TokenFindOptions{PlayerID: playerID},
		)
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return token, nil
	}
	return finded, nil
}

func (service *TokenService) IssueToken(playerID uint64) (*app_model.Token, error) {
	newJwt := jwt.IssueShort(
		&jwtmodels.JwtShortPayload{
			Iss:      common.AppName,
			Sub:      fmt.Sprint(playerID),
			Duration: time.Hour * 24 * 30 * 6,
		},
		common.Conf.SecretKey,
	)
	newTokenID, err := service.tokenRepository.Insert(
		&app_model.Token{PlayerID: playerID, Jwt: newJwt, Status: app_model.TokenAvailable},
	)
	if err != nil {
		return nil, err
	}
	newToken, err := service.tokenRepository.FindOne(&app_option.TokenFindOptions{ID: newTokenID})
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return newToken, err
}

func (service *TokenService) BlockToken(tokenID uint64) error {
	return service.tokenRepository.UpdateStatus(tokenID, app_model.TokenBlocked)
}
