package models

import adModels "api-tests-template/internal/managers/advertisements/models"

// описывает тело запроса добавления объявления в избранное
type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

// описывает запись в избранном
type Favorite struct {
	FavoriteID    string                 `json:"favorite_id"`
	AddedAt       string                 `json:"added_at"`
	Advertisement adModels.Advertisement `json:"advertisement"`
}

// описывает страницу избранного
type FavoritesList struct {
	Items      []Favorite `json:"items"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	Total      int        `json:"total"`
	TotalPages int        `json:"total_pages"`
}
