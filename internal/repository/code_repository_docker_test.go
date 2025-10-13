package repository_test

import (
	"fmt"
	"testing"

	"context"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/go-redis/redis/v8"
	"github.com/mbilarusdev/durak_auth_bot/internal/adapter"
	"github.com/mbilarusdev/durak_auth_bot/internal/repository"
	"github.com/mbilarusdev/durak_auth_bot/internal/utils"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCodeRepositoryDocker(t *testing.T) {
	rdbPass := "test_pass"
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image: "redis:alpine",
		Cmd: []string{
			"redis-server",
			"--requirepass",
			rdbPass,
			"--appendonly",
			"yes",
		},
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort(nat.Port("6379/tcp")),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Couldn't start container: %s", err)
	}
	defer container.Terminate(ctx)

	ip, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("Couldn't get host ip: %s", err)
	}

	port, err := container.MappedPort(ctx, "6379")
	if err != nil {
		t.Fatalf("Couldn't get mapped port: %s", err)
	}

	time.Sleep(time.Second * 5)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", ip, port.Port()),
		Password: rdbPass,
		DB:       0,
	})

	defer rdb.Close()

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		t.Fatal(err)
	}

	t.Run("All methods successs", func(t *testing.T) {
		codeRepository := repository.NewCodeRepository(adapter.NewCacheManagerAdapter(rdb))

		phoneNumber := "+79268566212"
		code := utils.GenerateRandomCode()

		err := codeRepository.Save(phoneNumber, code)
		if err != nil {
			t.Fatalf("Save code failed: %s", err)
		}

		findedCode, err := codeRepository.Get(phoneNumber)
		if err != nil {
			t.Fatalf("Get code failed: %s", err)
		}
		if findedCode != code {
			t.Fatalf("Finded code not equals to saved code")
		}

		err = codeRepository.Del(phoneNumber)
		if err != nil {
			t.Fatalf("Delete code failed: %s", err)
		}

		findedCode, err = codeRepository.Get(phoneNumber)
		if err == nil || err != redis.Nil {
			t.Fatalf("Expected redis.Nil error when get deleted code: %s", err)
		}
		if findedCode != "" {
			t.Fatalf("Finded code not equals to empty")
		}
	})
}
