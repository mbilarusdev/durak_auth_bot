package service

import (
	"fmt"

	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/jwt/jwt"
	jwtmodels "github.com/mbilarusdev/jwt/models"
)

type TokenManager interface {
	FindActualByToken(token string) (*models.Token, error)
	FindActualByPlayerID(playerID uint64) (*models.Token, error)
	IssueToken(playerID uint64) (*models.Token, error)
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

func (service *TokenService) FindActualByToken(token string) (*models.Token, error) {
	playerID := jwt.GetSubID(token, common.Conf.Token)
	return service.FindActualByPlayerID(playerID)
}

func (service *TokenService) FindActualByPlayerID(playerID uint64) (*models.Token, error) {

	finded, err := service.tokenRepository.FindOne(playerID)
	if err != nil {
		return nil, err
	}
	if finded == nil {
		return nil, nil
	}
	available := jwt.Check(finded.Jwt, common.Conf.SecretKey)
	if !available && finded.Status != models.TokenBlocked {
		err := service.tokenRepository.UpdateStatus(finded.ID, models.TokenExpired)
		if err != nil {
			return nil, err
		}
		return service.tokenRepository.FindOne(playerID)
	}
	return finded, nil
}

func (service *TokenService) IssueToken(playerID uint64) (*models.Token, error) {
	newJwt := jwt.IssueShort(
		&jwtmodels.JwtShortPayload{Iss: common.AppName, Sub: fmt.Sprint(playerID)},
		common.Conf.SecretKey,
	)
	newToken, err := service.tokenRepository.Insert(
		&models.Token{PlayerID: playerID, Jwt: newJwt, Status: models.TokenAvailable},
	)
	return newToken, err
}

func (service *TokenService) BlockToken(tokenID uint64) error {
	return service.tokenRepository.UpdateStatus(tokenID, models.TokenBlocked)
}
