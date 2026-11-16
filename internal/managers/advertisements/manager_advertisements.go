package advertisements

import (
	"net/http"
	"testing"

	advClient "api-tests-template/internal/client/http/advertisements"
	"api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/apierror"
	"api-tests-template/internal/utils"
)

// ищет объявления по подстроке в названии и возвращает список, ожидает 200
func GetAdvertisementsList(t *testing.T, token, text string) models.AdvertisementsList {
	return utils.DecodeBody[models.AdvertisementsList](t, advClient.HttpGetAdvertisements(t, token, text))
}

// создаёт объявление и возвращает его, ожидает 201
func CreateAdvertisement(t *testing.T, token string, req models.CreateAdvertisementRequest) models.Advertisement {
	return utils.DecodeBody[models.Advertisement](t, advClient.HttpPostAdvertisement(t, token, req, http.StatusCreated))
}

// пытается создать объявление и возвращает ошибку API с ожидаемым статус-кодом
func CreateAdvertisementExpectError(t *testing.T, token string, req models.CreateAdvertisementRequest, expectedStatus int) apierror.ErrorResponse {
	return utils.DecodeBody[apierror.ErrorResponse](t, advClient.HttpPostAdvertisement(t, token, req, expectedStatus))
}

// возвращает объявление по ID, ожидает 200
func GetAdvertisementByID(t *testing.T, token, id string) models.Advertisement {
	return utils.DecodeBody[models.Advertisement](t, advClient.HttpGetAdvertisementByID(t, token, id, http.StatusOK))
}

// возвращает фотографии объявления, ожидает 200
func GetAdvertisementPhotos(t *testing.T, token, id string) []models.Photo {
	return utils.DecodeBody[[]models.Photo](t, advClient.HttpGetAdvertisementPhotos(t, token, id, http.StatusOK))
}

// удаляет объявление по ID, ожидает 204
func DeleteAdvertisement(t *testing.T, token, id string) {
	resp := advClient.HttpDeleteAdvertisement(t, token, id, http.StatusNoContent)
	defer resp.Body.Close()
}
