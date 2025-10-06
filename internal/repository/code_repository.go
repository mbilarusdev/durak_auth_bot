package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mbilarusdev/durak_auth_bot/internal/locator"
)

type CodeProvider interface {
	SaveCode(phoneNumber string, code string) error
	GetCode(phoneNumber string) (string, error)
	DelCode(phoneNumber string) error
}

type CodeRepository struct{}

func NewCodeRepository() *CodeRepository {
	return new(CodeRepository)
}

func (repository *CodeRepository) SaveCode(phoneNumber string, code string) error {
	client := locator.Instance.Get("redis_client").(*redis.Client)
	_, err := client.Set(
		context.Background(),
		fmt.Sprintf("auth_bot:code:%v", phoneNumber),
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
	client := locator.Instance.Get("redis_client").(*redis.Client)
	code, err := client.Get(
		context.Background(),
		fmt.Sprintf("auth_bot:code:%v", phoneNumber),
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
	client := locator.Instance.Get("redis_client").(*redis.Client)
	_, err := client.Del(
		context.Background(),
		fmt.Sprintf("auth_bot:code:%v", phoneNumber),
	).Result()
	if err != nil {
		log.Println("Ошибка удаления кода из redis")
		return err
	}

	return nil
}
