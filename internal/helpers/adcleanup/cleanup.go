package adcleanup

import (
	"testing"

	"api-tests-template/internal/managers/advertisements"
)

// хранит объявления, созданные тестом, чтобы удалить их после его завершения
type Registry struct {
	tokensByAdID map[string]string
}

// создаёт пустой реестр
func New() *Registry {
	return &Registry{tokensByAdID: map[string]string{}}
}

// запоминает объявление и токен его владельца для последующего удаления
func (r *Registry) Track(ownerToken, adID string) {
	r.tokensByAdID[adID] = ownerToken
}

// убирает объявление из реестра (например, если тест уже удалил его сам)
func (r *Registry) Forget(adID string) {
	delete(r.tokensByAdID, adID)
}

// удаляет все запомненные объявления и очищает реестр
func (r *Registry) DeleteAll(t *testing.T) {
	for adID, ownerToken := range r.tokensByAdID {
		advertisements.DeleteAdvertisement(t, ownerToken, adID)
	}
	r.tokensByAdID = map[string]string{}
}
