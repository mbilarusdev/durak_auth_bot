package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
)

const (
	cacheKey string = "code"
)

type CodeProvider interface {
	SaveCode(phoneNumber string, code string) error
	GetCode(phoneNumber string) (string, error)
	DelCode(phoneNumber string) error
}

type CodeRepository struct {
	rdb interfaces.CacheManager
}

func NewCodeRepository(rdb *redis.Client) *CodeRepository {
	repository := new(CodeRepository)
	repository.rdb = rdb
	return repository
}

func (repository *CodeRepository) SaveCode(phoneNumber string, code string) error {
	_, err := repository.rdb.Set(
		context.Background(),
		fmt.Sprintf("%v:%v:%v", common.ServiceCacheName, cacheKey, phoneNumber),
		code,
		time.Minute,
	).Result()
	if err != nil {
		log.Println("Ошибка сохранения кода в redis")
		return err
	}
	return nil
}

func (repository *CodeRepository) GetCode(phoneNumber string) (string, error) {
	code, err := repository.rdb.Get(
		context.Background(),
		fmt.Sprintf("%v:%v:%v", common.ServiceCacheName, cacheKey, phoneNumber),
	).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		log.Println("Ошибка сохранения кода в redis")
		return "", err
	}

	return code, nil
}

func (repository *CodeRepository) DelCode(phoneNumber string) error {
	_, err := repository.rdb.Del(
		context.Background(),
		fmt.Sprintf("%v:%v:%v", common.ServiceCacheName, cacheKey, phoneNumber),
	).Result()
	if err != nil {
		log.Println("Ошибка удаления кода из redis")
		return err
	}

	return nil
}
