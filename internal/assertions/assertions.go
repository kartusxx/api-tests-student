package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/managers/advertisements/models"

	"api-tests-template/internal/managers/apierror"
)

func ContainsAdvertisementByTitle(t *testing.T, items []models.Advertisement, title string) {
	t.Helper()
	for _, item := range items {
		if item.Title == title {
			return
		}
	}
	require.Failf(t, "объявление не найдено", "в ответе нет объявления с названием %q", title)
}

// проверяет значения error и message в ответе API с ошибкой
func RequireAPIError(t *testing.T, actual apierror.ErrorResponse, expectedError, expectedMessage string) {
	t.Helper()
	require.Equal(t, expectedError, actual.Error, "Некорректное значение error параметра")
	require.Equal(t, expectedMessage, actual.Message, "Некорректное значение message параметра")
}
