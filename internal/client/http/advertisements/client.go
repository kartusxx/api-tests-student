package advertisements

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/advertisements/models"
)

// выполняет GET-запрос поиска объявлений по подстроке в названии, ожидает 200
func HttpGetAdvertisements(t *testing.T, token, textForSearch string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementsPath).
		Query("search", textForSearch).
		Expect(t).Status(http.StatusOK).End().Response
}

// создаёт объявление через multipart/form-data и ожидает указанный статус-код
func HttpPostAdvertisement(t *testing.T, token string, req models.CreateAdvertisementRequest, expectedStatus int) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", req.Title).
		MultipartFormData("description", req.Description)

	if req.Price != nil {
		request = request.MultipartFormData("price", strconv.FormatFloat(*req.Price, 'f', -1, 64))
	}
	if req.Quantity != nil {
		request = request.MultipartFormData("quantity", strconv.Itoa(*req.Quantity))
	}
	if len(req.PhotoPaths) > 0 {
		request = request.MultipartFile("photos", req.PhotoPaths...)
	}

	return request.Expect(t).Status(expectedStatus).End().Response
}

// получает объявление по ID и ожидает указанный статус-код
func HttpGetAdvertisementByID(t *testing.T, token, id string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).Status(expectedStatus).End().Response
}

// получает фотографии объявления и ожидает указанный статус-код
func HttpGetAdvertisementPhotos(t *testing.T, token, id string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Getf("%s/%s/photos", path.AdvertisementsPath, id).
		Expect(t).Status(expectedStatus).End().Response
}

// удаляет объявление по ID и ожидает указанный статус-код
func HttpDeleteAdvertisement(t *testing.T, token, id string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).Status(expectedStatus).End().Response
}
