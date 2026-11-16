package auth

import (
	"net/http"
	"testing"

	authClient "api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

const (
	testUserBirthDate = "2000-01-01T00:00:00Z"
	testUserFirstName = "Автотест"
	testUserLastName  = "Пользователь"
	testUserPassword  = "Test1234"
)

// регистрирует нового пользователя со случайным email и возвращает его токен и данные, ожидает 201
func RegisterRandomUser(t *testing.T) models.LoginOkResponse {
	request := models.RegisterRequest{
		BirthDate: testUserBirthDate,
		Email:     utils.RandomEmail(),
		FirstName: testUserFirstName,
		LastName:  testUserLastName,
		Password:  testUserPassword,
	}
	return utils.DecodeBody[models.LoginOkResponse](t, authClient.HttpPostAuthRegister(t, request, http.StatusCreated))
}
