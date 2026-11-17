package favorites

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/helpers/adcleanup"
	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/apierror"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	favManager "api-tests-template/internal/managers/favorites"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	adTitlePrefix = "Автотест избранное"
	adDescription = "Описание объявления для тестов избранного"
	photoFileName = "photo.png"

	nonexistentAdID = "11111111-2222-4333-8444-555555555555"
	wrongAuthValue  = "incorrect_token"

	errUnauthorized  = "unauthorized"
	errNotFound      = "not_found"
	errAlreadyExists = "already_exists"
	errForbidden     = "forbidden"

	msgAuthHeaderNeeded   = "Authorization header required"
	msgInvalidOrExpired   = "Invalid or expired token"
	msgAdNotFound         = "Advertisement not found"
	msgFavoriteNotFound   = "Favorite entry not found"
	msgAlreadyInFavorites = "Advertisement is already in favorites"
	msgOwnAdvertisement   = "Cannot add own advertisement to favorites"
)

type TestSuite struct {
	suite.Suite
	owner authModels.LoginOkResponse // владелец объявлений (отдельный зарегистрированный пользователь)
	user  authModels.LoginOkResponse // второй пользователь, который добавляет объявления в избранное
	ads   *adcleanup.Registry
	ad    adModels.Advertisement // объявление владельца, создаётся перед каждым тестом
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация владельца объявлений со случайным email")
	s.owner = auth.RegisterRandomUser(s.T())

	base.Precondition("Регистрация второго пользователя со случайным email")
	s.user = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) SetupTest() {
	base.Precondition("Владелец создаёт объявление для теста")
	s.ads = adcleanup.New()
	request := adModels.CreateAdvertisementRequest{
		Title:       utils.UniqueTitle(adTitlePrefix),
		Description: adDescription,
		PhotoPaths:  []string{utils.CreateTestPNG(s.T(), photoFileName)},
	}
	s.ad = advertisements.CreateAdvertisement(s.T(), s.owner.Token, request)
	s.ads.Track(s.owner.Token, s.ad.ID)
}

func (s *TestSuite) TearDownTest() {
	base.Precondition("Удаляем объявления, созданные в тесте")
	s.ads.DeleteAll(s.T())
}

// Positive: добавленное объявление появляется в избранном, а удалённое из избранного пропадает
func (s *TestSuite) TestFavoritesLifecycle() {
	s.Run("Добавляем объявление в избранное", func() {
		favorite := favManager.AddToFavorites(s.T(), s.user.Token, s.ad.ID)
		require.Equal(s.T(), s.ad.ID, favorite.Advertisement.ID, "В ответе добавления другое объявление")
	})

	s.Run("Объявление появилось в списке избранного", func() {
		list := favManager.GetFavorites(s.T(), s.user.Token, s.ad.Title)
		assertions.ContainsFavoriteByAdID(s.T(), list.Items, s.ad.ID)
	})

	s.Run("Удаляем объявление из избранного", func() {
		favManager.RemoveFromFavorites(s.T(), s.user.Token, s.ad.ID)
	})

	s.Run("Объявление пропало из списка избранного", func() {
		list := favManager.GetFavorites(s.T(), s.user.Token, s.ad.Title)
		assertions.NotContainsFavoriteByAdID(s.T(), list.Items, s.ad.ID)
	})
}

// Positive: если владелец удаляет само объявление, оно исчезает из избранного пользователя
func (s *TestSuite) TestFavoritesAdvertisementDeleted() {
	s.Run("Добавляем объявление в избранное", func() {
		favManager.AddToFavorites(s.T(), s.user.Token, s.ad.ID)
		list := favManager.GetFavorites(s.T(), s.user.Token, s.ad.Title)
		assertions.ContainsFavoriteByAdID(s.T(), list.Items, s.ad.ID)
	})

	s.Run("Владелец удаляет объявление", func() {
		advertisements.DeleteAdvertisement(s.T(), s.owner.Token, s.ad.ID)
		s.ads.Forget(s.ad.ID)
	})

	s.Run("Объявление пропало из избранного пользователя", func() {
		list := favManager.GetFavorites(s.T(), s.user.Token, s.ad.Title)
		assertions.NotContainsFavoriteByAdID(s.T(), list.Items, s.ad.ID)
	})
}

// Negative: операции с несуществующим объявлением возвращают 404
func (s *TestSuite) TestFavoritesNonexistentAdvertisement() {
	s.Run("Добавление несуществующего объявления в избранное", func() {
		apiError := favManager.AddToFavoritesExpectError(s.T(), s.user.Token, nonexistentAdID, http.StatusNotFound)
		assertions.RequireAPIError(s.T(), apiError, errNotFound, msgAdNotFound)
	})

	s.Run("Удаление несуществующего объявления из избранного", func() {
		apiError := favManager.RemoveFromFavoritesExpectError(s.T(), s.user.Token, nonexistentAdID, http.StatusNotFound)
		assertions.RequireAPIError(s.T(), apiError, errNotFound, msgFavoriteNotFound)
	})
}

// Negative: ручки избранного без авторизации или с неверным токеном возвращают 401
func (s *TestSuite) TestFavoritesUnauthorized() {
	operations := []struct {
		name string
		call func(token string) apierror.ErrorResponse
	}{
		{"получение избранного", func(token string) apierror.ErrorResponse {
			return favManager.GetFavoritesExpectError(s.T(), token, http.StatusUnauthorized)
		}},
		{"добавление в избранное", func(token string) apierror.ErrorResponse {
			return favManager.AddToFavoritesExpectError(s.T(), token, s.ad.ID, http.StatusUnauthorized)
		}},
		{"удаление из избранного", func(token string) apierror.ErrorResponse {
			return favManager.RemoveFromFavoritesExpectError(s.T(), token, s.ad.ID, http.StatusUnauthorized)
		}},
	}

	authCases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", msgAuthHeaderNeeded},
		{"неверный токен", wrongAuthValue, msgInvalidOrExpired},
	}

	for _, authCase := range authCases {
		for _, operation := range operations {
			s.Run(authCase.name+": "+operation.name, func() {
				apiError := operation.call(authCase.token)
				assertions.RequireAPIError(s.T(), apiError, errUnauthorized, authCase.expectedMessage)
			})
		}
	}
}

// Negative: повторное добавление того же объявления возвращает 409
func (s *TestSuite) TestFavoritesAddTwice() {
	s.Run("Добавляем объявление в избранное первый раз", func() {
		favManager.AddToFavorites(s.T(), s.user.Token, s.ad.ID)
	})

	s.Run("Повторное добавление возвращает конфликт", func() {
		apiError := favManager.AddToFavoritesExpectError(s.T(), s.user.Token, s.ad.ID, http.StatusConflict)
		assertions.RequireAPIError(s.T(), apiError, errAlreadyExists, msgAlreadyInFavorites)
	})
}

// Negative: нельзя добавить в избранное своё собственное объявление (403)
func (s *TestSuite) TestFavoritesAddOwnAdvertisement() {
	s.Run("Владелец добавляет своё объявление в избранное", func() {
		apiError := favManager.AddToFavoritesExpectError(s.T(), s.owner.Token, s.ad.ID, http.StatusForbidden)
		assertions.RequireAPIError(s.T(), apiError, errForbidden, msgOwnAdvertisement)
	})
}
