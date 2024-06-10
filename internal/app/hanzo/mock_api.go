package hanzo

import (
	"fmt"
	"hanzotg/internal/app/utils"
	"log"
)

type MockApiClient struct{}

func (m *MockApiClient) AddUser(userId, username string) error {
	log.Println(fmt.Sprintf("Adding user [%s] %s to hanzo db", userId, username))
	return nil
}

func (m *MockApiClient) GenerateUserPassword(userId string) (string, error) {
	log.Println(fmt.Sprintf("Generating password for user [%s]", userId))
	return utils.RandomString(16), nil
}

func (m *MockApiClient) DeleteUser(userId string) error {
	log.Println(fmt.Sprintf("Deleting password for [%s]", userId))
	return nil
}
