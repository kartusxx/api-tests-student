package createAdvertisement

import (
	"os"
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

	base.Precondition("Авторизация пользователя с кредами из переменных окружения и получение его параметров")
	s.loginData = auth.Login(s.T(), os.Getenv("TEST_LOGIN"), os.Getenv("TEST_PASSWORD"))
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
