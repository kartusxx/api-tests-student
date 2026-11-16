package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
	"api-tests-template/internal/managers/favorites/models"
)

// добавляет объявление в избранное и ожидает указанный статус-код
func HttpPostFavorite(t *testing.T, token, adID string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		JSON(models.AddFavoriteRequest{AdvertisementID: adID}).
		Expect(t).Status(expectedStatus).End().Response
}

// получает список избранного (с необязательным поиском по названию) и ожидает указанный статус-код
func HttpGetFavorites(t *testing.T, token, search string, expectedStatus int) *http.Response {
	request := apiRunner.GetRunner().Auth(token).Create().Get(path.FavoritesPath)
	if search != "" {
		request = request.Query("search", search)
	}
	return request.Expect(t).Status(expectedStatus).End().Response
}

// удаляет объявление из избранного и ожидает указанный статус-код
func HttpDeleteFavorite(t *testing.T, token, adID string, expectedStatus int) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Deletef("%s/%s", path.FavoritesPath, adID).
		Expect(t).Status(expectedStatus).End().Response
}
