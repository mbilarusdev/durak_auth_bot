package repository_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/interfaces"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/utils"
	"github.com/stretchr/testify/mock"
)

func TestCodeRepository(t *testing.T) {
	phoneNumber := "+79150719588"

	t.Run("Code Save with success", func(t *testing.T) {
		code := utils.GenerateRandomCode()
		rdb := interfaces.NewMockCacheManager(t)
		statusCmd := interfaces.NewMockCacheStatusCmd(t)
		statusCmd.On("Result").Return("ok", nil)
		rdb.On(
			"Set",
			mock.Anything,
			fmt.Sprintf(
				"%v:%v:%v",
				common.ServiceCacheName,
				repository.CodeRepoCacheKey,
				phoneNumber,
			),
			code,
			time.Minute,
		).Return(statusCmd)

		codeRepository := repository.NewCodeRepository(rdb)
		err := codeRepository.Save(phoneNumber, code)

		if err != nil {
			t.Errorf("Save() failed: %v", err)
		}
	})

	t.Run("Code Get with success", func(t *testing.T) {
		code := utils.GenerateRandomCode()
		rdb := interfaces.NewMockCacheManager(t)
		stringCmd := interfaces.NewMockCacheStringCmd(t)
		stringCmd.On("Result").Return(code, nil)
		rdb.On(
			"Get",
			mock.Anything,
			fmt.Sprintf(
				"%v:%v:%v",
				common.ServiceCacheName,
				repository.CodeRepoCacheKey,
				phoneNumber,
			),
		).Return(stringCmd)

		codeRepository := repository.NewCodeRepository(rdb)
		res, err := codeRepository.Get(phoneNumber)

		if res != code && err != nil {
			t.Errorf("Get() failed: %v", err)
		}
	})

	t.Run("Code del with success", func(t *testing.T) {
		rdb := interfaces.NewMockCacheManager(t)
		intCmd := interfaces.NewMockCacheIntCmd(t)
		intCmd.On("Result").Return(int64(1), nil)
		rdb.On(
			"Del",
			mock.Anything,
			[]string{fmt.Sprintf(
				"%v:%v:%v",
				common.ServiceCacheName,
				repository.CodeRepoCacheKey,
				phoneNumber,
			)},
		).Return(intCmd)

		codeRepository := repository.NewCodeRepository(rdb)
		err := codeRepository.Del(phoneNumber)

		if err != nil {
			t.Errorf("Del() failed: %v", err)
		}
	})

	t.Run("Code get with not finded", func(t *testing.T) {
		rdb := interfaces.NewMockCacheManager(t)
		stringCmd := interfaces.NewMockCacheStringCmd(t)
		stringCmd.On("Result").Return("", redis.Nil)
		rdb.On(
			"Get",
			mock.Anything,
			fmt.Sprintf(
				"%v:%v:%v",
				common.ServiceCacheName,
				repository.CodeRepoCacheKey,
				phoneNumber,
			),
		).Return(stringCmd)

		codeRepository := repository.NewCodeRepository(rdb)
		res, err := codeRepository.Get(phoneNumber)

		if res != "" || err == nil || err != redis.Nil {
			t.Errorf("Get() expected redis.Nil err but return nil or other error")
		}
	})

	t.Run("Code del with not finded", func(t *testing.T) {
		rdb := interfaces.NewMockCacheManager(t)
		intCmd := interfaces.NewMockCacheIntCmd(t)
		intCmd.On("Result").Return(int64(0), redis.Nil)
		rdb.On(
			"Del",
			mock.Anything,
			[]string{fmt.Sprintf(
				"%v:%v:%v",
				common.ServiceCacheName,
				repository.CodeRepoCacheKey,
				phoneNumber,
			)},
		).Return(intCmd)

		codeRepository := repository.NewCodeRepository(rdb)
		err := codeRepository.Del(phoneNumber)

		if err == nil || err != redis.Nil {
			t.Errorf("Del() expected redis.Nil err but return nil or other error")
		}
	})
}
