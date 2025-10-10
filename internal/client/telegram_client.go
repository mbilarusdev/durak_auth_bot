package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mbilarusdev/durak_auth_bot/internal/common"
	"github.com/mbilarusdev/durak_auth_bot/internal/models"
)

const (
	tgApiUrl string = "https://api.telegram.org/bot"
)

type TgNetworkManager interface {
	Get(method string) ([]byte, error)
	Post(method string, data []byte) (*models.SendResponse, error)
}

type TelegramClient struct{}

func NewTelegramClient() *TelegramClient {
	return new(TelegramClient)
}

func (tgClient *TelegramClient) Get(method string) ([]byte, error) {

	url := fmt.Sprintf(tgApiUrl+"%s/%s", common.Conf.Token, method)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body, nil
}

func (tgClient *TelegramClient) Post(method string, data []byte) (*models.SendResponse, error) {
	url := fmt.Sprintf(tgApiUrl+"%s/%s", common.Conf.Token, method)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	responseBody, _ := io.ReadAll(resp.Body)
	var sendResp models.SendResponse
	json.Unmarshal(responseBody, &sendResp)
	return &sendResp, nil
}
