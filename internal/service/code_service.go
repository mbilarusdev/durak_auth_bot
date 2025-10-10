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
	return smsCode, service.codeRepository.SaveCode(phone, smsCode)
}

func (service *CodeService) ConsumeIsRightCode(phone string, confirmCode string) (bool, error) {
	code, err := service.codeRepository.GetCode(phone)
	if err != nil {
		return false, err
	}
	isRightCode := code == confirmCode
	if isRightCode {
		service.codeRepository.DelCode(phone)
	}
	return isRightCode, nil
}
