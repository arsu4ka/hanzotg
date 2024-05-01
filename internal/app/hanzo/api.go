package hanzo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hanzotg/internal/app/utils"
	"net/http"
	"strings"
)

type ApiConfig struct {
	SecretKey string
	BaseUrl   string
}

type ApiClient struct {
	client *http.Client
	ApiConfig
}

func NewApiClient(config ApiConfig) *ApiClient {
	config.BaseUrl = strings.TrimSuffix(config.BaseUrl, "/")
	return &ApiClient{
		client:    http.DefaultClient,
		ApiConfig: config,
	}
}

func (a *ApiClient) AddUser(userId, username string) error {
	url := fmt.Sprintf(
		"%s/add_user?userId=%s&username=%s&chatId=%s",
		a.BaseUrl,
		userId,
		username,
		userId,
	)
	_, err := a.client.Get(url)
	return err
}

func (a *ApiClient) GenerateUserPassword(userId string) (string, error) {
	var apiResponse struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	url := a.BaseUrl + "/add_user_with_chat_id/"
	body := map[string]string{
		"chat_id":    userId,
		"secret_key": a.SecretKey,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	response, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer utils.SafeCloseBody(response)

	if err := json.NewDecoder(response.Body).Decode(&apiResponse); err != nil {
		return "", err
	}

	return apiResponse.Password, nil
}

func (a *ApiClient) DeleteUser(userId string) error {
	url := a.BaseUrl + "/delete_user_with_chat_id/"
	body := map[string]string{
		"chat_id":    userId,
		"secret_key": a.SecretKey,
	}

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = a.client.Do(req)
	return err
}
