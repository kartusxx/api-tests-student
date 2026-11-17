package assertions

import (
	"testing"

	"github.com/stretchr/testify/require"

	favModels "api-tests-template/internal/managers/favorites/models"
)

// проверяет, что в списке избранного есть объявление с указанным id
func ContainsFavoriteByAdID(t *testing.T, items []favModels.Favorite, adID string) {
	t.Helper()
	for _, item := range items {
		if item.Advertisement.ID == adID {
			return
		}
	}
	require.Failf(t, "объявление не найдено в избранном", "в избранном нет объявления с id %q", adID)
}

// проверяет, что в списке избранного нет объявления с указанным id
func NotContainsFavoriteByAdID(t *testing.T, items []favModels.Favorite, adID string) {
	t.Helper()
	for _, item := range items {
		if item.Advertisement.ID == adID {
			require.Failf(t, "объявление осталось в избранном", "в избранном присутствует объявление с id %q", adID)
		}
	}
}
