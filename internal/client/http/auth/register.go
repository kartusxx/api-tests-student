package auth

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/auth/models"
)

// регистрирует нового пользователя и ожидает указанный статус-код
func HttpPostAuthRegister(t *testing.T, request models.RegisterRequest, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Create().
		Post(path.AuthRegisterPath).
		JSON(request).
		Expect(t).Status(expectedStatus).End().Response
}
