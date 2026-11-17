package createAdvertisement

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"api-tests-template/internal/assertions"
	"api-tests-template/internal/helpers/adcleanup"
	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	authModels "api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	adTitlePrefix = "Автотест объявление"
	adDescription = "Описание тестового объявления"
	adPrice       = 1500.5
	adQuantity    = 3
	photoFileName = "photo.png"

	maxTitleLength = 50
	maxPhotosCount = 3
	negativePrice  = -5.0

	errInvalidRequest = "invalid_request"
	errUnauthorized   = "unauthorized"

	msgTitleRequired    = "title is required"
	msgPriceOutOfRange  = "price must be between 0 and 1000000 inclusive"
	msgPhotoRequired    = "at least 1 photo is required"
	msgTooManyPhotos    = "maximum 3 photos allowed"
	msgAuthHeaderNeeded = "Authorization header required"
	msgInvalidOrExpired = "Invalid or expired token"

	wrongAuthValue = "incorrect_token"
)

type TestSuite struct {
	suite.Suite
	loginData authModels.LoginOkResponse
	ads       *adcleanup.Registry
	photoPath string
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация отдельного пользователя для тестов со случайным email")
	s.loginData = auth.RegisterRandomUser(s.T())
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) SetupTest() {
	base.Precondition("Готовим тестовую фотографию и пустой реестр созданных объявлений")
	s.ads = adcleanup.New()
	s.photoPath = utils.CreateTestPNG(s.T(), photoFileName)
}

func (s *TestSuite) TearDownTest() {
	base.Precondition("Удаляем объявления, созданные в тесте")
	s.ads.DeleteAll(s.T())
}

// Positive: создание объявления со всеми параметрами и проверка через разные ручки
func (s *TestSuite) TestCreateAdvertisementPositive() {
	title := utils.UniqueTitle(adTitlePrefix)
	request := adModels.CreateAdvertisementRequest{
		Title:       title,
		Description: adDescription,
		Price:       utils.Ptr(adPrice),
		Quantity:    utils.Ptr(adQuantity),
		PhotoPaths:  []string{s.photoPath},
	}

	var created adModels.Advertisement
	s.Run("Создаём объявление со всеми параметрами", func() {
		created = advertisements.CreateAdvertisement(s.T(), s.loginData.Token, request)
		s.ads.Track(s.loginData.Token, created.ID)
		require.NotEmpty(s.T(), created.ID, "В ответе на создание отсутствует id объявления")
	})

	s.Run("Объявление доступно по ID и содержит переданные параметры", func() {
		got := advertisements.GetAdvertisementByID(s.T(), s.loginData.Token, created.ID)
		require.Equal(s.T(), created.ID, got.ID, "id объявления не совпадает")
		require.Equal(s.T(), title, got.Title, "title не совпадает")
		require.Equal(s.T(), adDescription, got.Description, "description не совпадает")
		require.InDelta(s.T(), adPrice, got.Price, 0.001, "price не совпадает")
		require.Equal(s.T(), adQuantity, got.Quantity, "quantity не совпадает")
		require.Equal(s.T(), s.loginData.User.Id, got.UserID, "объявление не принадлежит нашему пользователю")
	})

	s.Run("Объявление находится поиском по полному названию", func() {
		list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, title)
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, title)
	})

	s.Run("Фотография загружена и доступна по ссылке", func() {
		photos := advertisements.GetAdvertisementPhotos(s.T(), s.loginData.Token, created.ID)
		require.Len(s.T(), photos, 1, "Количество загруженных фотографий не совпадает")
		require.NotEmpty(s.T(), photos[0].URL, "У фотографии пустой url")
	})
}

// Negative: создание объявления с некорректными параметрами возвращает 400 и понятное сообщение
func (s *TestSuite) TestCreateAdvertisementInvalidParams() {
	validRequest := func() adModels.CreateAdvertisementRequest {
		return adModels.CreateAdvertisementRequest{
			Title:       utils.UniqueTitle(adTitlePrefix),
			Description: adDescription,
			PhotoPaths:  []string{s.photoPath},
		}
	}

	cases := []struct {
		name            string
		modify          func(req *adModels.CreateAdvertisementRequest)
		expectedMessage string
	}{
		{
			name:            "пустой title",
			modify:          func(req *adModels.CreateAdvertisementRequest) { req.Title = "" },
			expectedMessage: msgTitleRequired,
		},
		{
			name:            "отрицательная цена",
			modify:          func(req *adModels.CreateAdvertisementRequest) { req.Price = utils.Ptr(negativePrice) },
			expectedMessage: msgPriceOutOfRange,
		},
		{
			name: "слишком длинный title",
			modify: func(req *adModels.CreateAdvertisementRequest) {
				req.Title = strings.Repeat("а", maxTitleLength+1)
			},
			expectedMessage: fmt.Sprintf("title must not exceed %d characters provided %d", maxTitleLength, maxTitleLength+1),
		},
		{
			name:            "нет фотографий",
			modify:          func(req *adModels.CreateAdvertisementRequest) { req.PhotoPaths = nil },
			expectedMessage: msgPhotoRequired,
		},
		{
			name: "больше трёх фотографий",
			modify: func(req *adModels.CreateAdvertisementRequest) {
				req.PhotoPaths = make([]string, maxPhotosCount+1)
				for i := range req.PhotoPaths {
					req.PhotoPaths[i] = s.photoPath
				}
			},
			expectedMessage: msgTooManyPhotos,
		},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			request := validRequest()
			tc.modify(&request)

			apiError := advertisements.CreateAdvertisementExpectError(s.T(), s.loginData.Token, request, http.StatusBadRequest)
			assertions.RequireAPIError(s.T(), apiError, errInvalidRequest, tc.expectedMessage)
		})
	}
}

// Negative: создание объявления без токена или с неверным токеном возвращает 401
func (s *TestSuite) TestCreateAdvertisementUnauthorized() {
	request := adModels.CreateAdvertisementRequest{
		Title:       utils.UniqueTitle(adTitlePrefix),
		Description: adDescription,
		PhotoPaths:  []string{s.photoPath},
	}

	cases := []struct {
		name            string
		token           string
		expectedMessage string
	}{
		{"без токена", "", msgAuthHeaderNeeded},
		{"неверный токен", wrongAuthValue, msgInvalidOrExpired},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			apiError := advertisements.CreateAdvertisementExpectError(s.T(), tc.token, request, http.StatusUnauthorized)
			assertions.RequireAPIError(s.T(), apiError, errUnauthorized, tc.expectedMessage)
		})
	}
}
