package myAdvertisement

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	"api-tests-template/internal/helpers/adcleanup"
	"api-tests-template/internal/managers/advertisements"
	adModels "api-tests-template/internal/managers/advertisements/models"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/managers/myAdvertisements"
	"api-tests-template/internal/utils"

	base "api-tests-template/tests"
)

const (
	adTitlePrefix  = "Автотест мои объявления"
	adDescription  = "Описание объявления для теста собственных объявлений"
	photoFileName  = "photo.png"
	myAdsPageLimit = 50

	errUnauthorized     = "unauthorized"
	msgInvalidOrExpired = "Invalid or expired token"
	wrongAuthValue      = "incorrect_token"
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
	ads       *adcleanup.Registry
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация отдельного пользователя со случайным email")
	s.loginData = auth.RegisterRandomUser(s.T())

	base.Precondition("Создаём объявление, которое должно появиться в списке собственных объявлений")
	s.ads = adcleanup.New()
	created := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, adModels.CreateAdvertisementRequest{
		Title:       utils.UniqueTitle(adTitlePrefix),
		Description: adDescription,
		PhotoPaths:  []string{utils.CreateTestPNG(s.T(), photoFileName)},
	})
	s.ads.Track(s.loginData.Token, created.ID)
}

func (s *TestSuite) TearDownSuite() {
	base.Precondition("Удаляем объявление, созданное для теста")
	s.ads.DeleteAll(s.T())
	base.TearDownSuite()
}

// Positive: Проверяем получение собственных объявлений в профиле
func (s *TestSuite) TestGetMyAdvertisementsPositive() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), s.loginData.Token, http.StatusOK, myAdsPageLimit)
	})

	var items gjson.Result
	s.Run("Проверяем, что у нас есть добавленное ранее объявление", func() {
		items = gjson.Get(advertisementsBody, "items")
		require.NotEmpty(s.T(), items.Array())
	})

	s.Run("Проверяем, что наши объявления имеют принадлежность к нашему пользователю и объявления имеют НЕ пустые параметры", func() {
		for _, item := range items.Array() {
			require.Equal(s.T(), s.loginData.User.Id, gjson.Get(item.String(), "user_id").String(), "user_id не совпадает")

			require.NotEmpty(s.T(), gjson.Get(item.String(), "id").String(), "id пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "created_at").String(), "created_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "updated_at").String(), "updated_at пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "title").String(), "Название пустое")
			require.NotEmpty(s.T(), gjson.Get(item.String(), "description").String(), "Описание пустое")
			require.True(s.T(), len(gjson.Get(item.String(), "photos").Array()) > 0, "Присутствуют пустые фотографии")
		}
	})
}

// Negative: Проверяем 401 ошибку при попытке получить собственные объявления с неверным токеном
func (s *TestSuite) TestGetMyAdvertisementsIncorrectToken() {
	var advertisementsBody string
	s.Run("Получаем список собственных объявлений с неправильным авторизационным токеном", func() {
		advertisementsBody = myAdvertisements.GetMyAdvertisements(s.T(), wrongAuthValue, http.StatusUnauthorized, myAdsPageLimit)
	})

	s.Run("Проверяем, что показывается правильная ошибка", func() {
		require.Equal(s.T(), errUnauthorized, gjson.Get(advertisementsBody, "error").String(),
			"Некорректное значение error параметра")
		require.Equal(s.T(), msgInvalidOrExpired, gjson.Get(advertisementsBody, "message").String(),
			"Некорректное значение message параметра")
	})
}
