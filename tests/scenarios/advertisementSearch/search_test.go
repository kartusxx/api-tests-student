package advertisementSearch

import (
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
	adTitleBase       = "Велосипед Горный Ромашка"
	adTitleStartPart  = "Велосипед Гор"
	adTitleMiddlePart = "Горный Ром"
	adTitleEndTrim    = "Велосипед Горный "
	adDescription     = "Описание объявления для тестов поиска"
	photoFileName     = "photo.png"
)

type TestSuite struct {
	suite.Suite
	loginData       authModels.LoginOkResponse
	ads             *adcleanup.Registry
	existingAdTitle string // полное название объявления, созданного для тестов
	searchBeginning string // подстрока из начала названия
	searchMiddle    string // подстрока из середины названия
	searchEnd       string // подстрока из конца названия
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Регистрация отдельного пользователя со случайным email")
	s.loginData = auth.RegisterRandomUser(s.T())

	base.Precondition("Создаём объявление, которое будем искать")
	s.ads = adcleanup.New()
	s.existingAdTitle = utils.UniqueTitle(adTitleBase)
	s.searchBeginning = adTitleStartPart
	s.searchMiddle = adTitleMiddlePart
	s.searchEnd = strings.TrimPrefix(s.existingAdTitle, adTitleEndTrim)

	created := advertisements.CreateAdvertisement(s.T(), s.loginData.Token, adModels.CreateAdvertisementRequest{
		Title:       s.existingAdTitle,
		Description: adDescription,
		PhotoPaths:  []string{utils.CreateTestPNG(s.T(), photoFileName)},
	})
	s.ads.Track(s.loginData.Token, created.ID)
}

func (s *TestSuite) TearDownSuite() {
	base.Precondition("Удаляем объявление, созданное для тестов поиска")
	s.ads.DeleteAll(s.T())
	base.TearDownSuite()
}

// Positive: поиск объявления по подстроке в названии
func (s *TestSuite) TestSearchByTitle() {
	list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, s.searchMiddle)

	s.Run("В ответе есть хотя бы одно объявление", func() {
		require.NotEmpty(s.T(), list.Items)
	})

	s.Run("Среди найденных есть наше объявление", func() {
		assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.existingAdTitle)
	})
}

// Positive: подстрока в начале, середине, конце названия и в разных регистрах
func (s *TestSuite) TestSearchByTitle_Substring() {
	cases := []struct {
		name   string
		search string
	}{
		{"начало", s.searchBeginning},
		{"середина", s.searchMiddle},
		{"конец", s.searchEnd},
		{"нижний регистр", strings.ToLower(s.searchMiddle)},
		{"верхний регистр", strings.ToUpper(s.searchMiddle)},
		{"смешанный регистр", utils.MixCase(s.searchMiddle)},
	}

	for _, tc := range cases {
		s.Run(tc.name, func() {
			list := advertisements.GetAdvertisementsList(s.T(), s.loginData.Token, tc.search)
			assertions.ContainsAdvertisementByTitle(s.T(), list.Items, s.existingAdTitle)
		})
	}
}
