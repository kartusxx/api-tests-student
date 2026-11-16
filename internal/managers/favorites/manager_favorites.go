package favorites

import (
	"net/http"
	"testing"

	favoritesClient "api-tests-template/internal/client/http/favorites"
	"api-tests-template/internal/managers/apierror"
	"api-tests-template/internal/managers/favorites/models"
	"api-tests-template/internal/utils"
)

// добавляет объявление в избранное и возвращает запись, ожидает 201
func AddToFavorites(t *testing.T, token, adID string) models.Favorite {
	return utils.DecodeBody[models.Favorite](t, favoritesClient.HttpPostFavorite(t, token, adID, http.StatusCreated))
}

// пытается добавить объявление в избранное и возвращает ошибку API с ожидаемым статус-кодом
func AddToFavoritesExpectError(t *testing.T, token, adID string, expectedStatus int) apierror.ErrorResponse {
	return utils.DecodeBody[apierror.ErrorResponse](t, favoritesClient.HttpPostFavorite(t, token, adID, expectedStatus))
}

// возвращает избранное пользователя (с необязательным поиском по названию), ожидает 200
func GetFavorites(t *testing.T, token, search string) models.FavoritesList {
	return utils.DecodeBody[models.FavoritesList](t, favoritesClient.HttpGetFavorites(t, token, search, http.StatusOK))
}

// пытается получить избранное и возвращает ошибку API с ожидаемым статус-кодом
func GetFavoritesExpectError(t *testing.T, token string, expectedStatus int) apierror.ErrorResponse {
	return utils.DecodeBody[apierror.ErrorResponse](t, favoritesClient.HttpGetFavorites(t, token, "", expectedStatus))
}

// удаляет объявление из избранного, ожидает 204
func RemoveFromFavorites(t *testing.T, token, adID string) {
	resp := favoritesClient.HttpDeleteFavorite(t, token, adID, http.StatusNoContent)
	defer resp.Body.Close()
}

// пытается удалить объявление из избранного и возвращает ошибку API с ожидаемым статус-кодом
func RemoveFromFavoritesExpectError(t *testing.T, token, adID string, expectedStatus int) apierror.ErrorResponse {
	return utils.DecodeBody[apierror.ErrorResponse](t, favoritesClient.HttpDeleteFavorite(t, token, adID, expectedStatus))
}
