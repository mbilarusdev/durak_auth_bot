package service

import (
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/utils"
)

type CodeManager interface {
	CreateCode(phone string) (string, error)
	ConsumeIsRightCode(phone string, confirmCode string) (bool, error)
}

type CodeService struct {
	codeRepository repository.CodeProvider
}

func NewCodeService(codeRepository repository.CodeProvider) *CodeService {
	service := new(CodeService)
	service.codeRepository = codeRepository
	return service
}

func (service *CodeService) CreateCode(phone string) (string, error) {
	smsCode := utils.GenerateRandomCode()
	err := service.codeRepository.Save(phone, smsCode)
	if err != nil {
		return "", err
	}
	return smsCode, nil
}

func (service *CodeService) ConsumeIsRightCode(phone string, confirmCode string) (bool, error) {
	code, err := service.codeRepository.Get(phone)
	if err != nil {
		return false, err
	}
	isRightCode := code == confirmCode
	if isRightCode {
		err := service.codeRepository.Del(phone)
		if err != nil {
			return isRightCode, err
		}
	}
	return isRightCode, nil
}
