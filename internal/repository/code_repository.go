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
	CodeRepoCacheKey string = "code"
)

type CodeProvider interface {
	Save(phoneNumber string, code string) error
	Get(phoneNumber string) (string, error)
	Del(phoneNumber string) error
}

type CodeRepository struct {
	rdb interfaces.CacheManager
}

func NewCodeRepository(rdb interfaces.CacheManager) *CodeRepository {
	repository := new(CodeRepository)
	repository.rdb = rdb
	return repository
}

func (repository *CodeRepository) Save(phoneNumber string, code string) error {
	_, err := repository.rdb.Set(
		context.Background(),
		fmt.Sprintf("%v:%v:%v", common.ServiceCacheName, CodeRepoCacheKey, phoneNumber),
		code,
		time.Minute,
	).Result()
	if err != nil {
		log.Println("Ошибка сохранения кода в redis")
		return err
	}
	return nil
}

func (repository *CodeRepository) Get(phoneNumber string) (string, error) {
	code, err := repository.rdb.Get(
		context.Background(),
		fmt.Sprintf("%v:%v:%v", common.ServiceCacheName, CodeRepoCacheKey, phoneNumber),
	).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Не найдено токена с номером телефона %v", phoneNumber)
			return "", err
		}
		log.Println("Ошибка сохранения кода в redis")
		return "", err
	}

	return code, nil
}

func (repository *CodeRepository) Del(phoneNumber string) error {
	_, err := repository.rdb.Del(
		context.Background(),
		fmt.Sprintf("%v:%v:%v", common.ServiceCacheName, CodeRepoCacheKey, phoneNumber),
	).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Не найдено токена с номером телефона %v", phoneNumber)
			return err
		}
		log.Println("Ошибка удаления кода из redis")
		return err
	}

	return nil
}
