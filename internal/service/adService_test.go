package service

import (
	"context"
	"github.com/AirstaNs/ValidationAds"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/entities"
	mocks "homework10/internal/mocks/repomocks"
	"homework10/internal/util"
	"strings"
	"testing"
	"time"
)

type serviceSuite struct {
	suite.Suite
	service   AdService
	formatter util.DateTimeFormatter
	adRepo    *mocks.AdRepository
	util.UID
}

var (
	wrongEmptyStr = ""
	wrongMoreStr  = strings.Repeat("r", 501)
	badID         = int64(-4124)

	testID = int64(0)

	testAd = entities.Ad{
		Title:     "Test",
		Text:      "TestText",
		AuthorID:  int64(1),
		Published: false,
	}

	dFilters = AdFilters{
		AuthorID:   int64(-1),
		CreateDate: time.Time{},
		Title:      "",
		Published:  true,
	}
)

func TestSuiteAdService(t *testing.T) {
	u := new(serviceSuite)
	suite.Run(t, u)

}

func (s *serviceSuite) SetupSuite() {
	AdRepo := new(mocks.AdRepository)
	formatter := util.NewDateTimeFormatter(time.DateOnly)
	s.service = NewAdsService(AdRepo, formatter)
	s.formatter = formatter
	s.adRepo = AdRepo

	toTime, err2 := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err2)

	testAd.CreateDate = toTime
	testAd.UpdateDate = toTime
	testAd.ID = testID

	s.adRepo.
		On("GetAdByID", testID).
		Return(&testAd, nil)

	s.adRepo.
		On("AddAd", testAd).
		Return(testID, nil)

}

func (s *serviceSuite) TearDownSuite() {
	s.service = nil
}

func (s *serviceSuite) Test_AdService_CreateAd() {

	ad, err := s.service.CreateAd(context.Background(), testAd.Title, testAd.Text, testAd.AuthorID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), *ad, testAd)
}

func (s *serviceSuite) Test_AdService_CreateAd_WrongTitle() {
	cAd := testAd
	cAd.Title = wrongEmptyStr

	_, err := s.service.CreateAd(context.Background(), cAd.Title, cAd.Text, cAd.AuthorID)
	assert.ErrorIs(s.T(), ValidationAds.ErrBadTitle, err)
}

func (s *serviceSuite) Test_AdService_CreateAd_WrongTitleMoreLen() {
	cAd := testAd
	cAd.Title = wrongMoreStr

	_, err := s.service.CreateAd(context.Background(), cAd.Title, cAd.Text, cAd.AuthorID)
	assert.ErrorIs(s.T(), ValidationAds.ErrBadTitle, err)
}

func (s *serviceSuite) Test_AdService_CreateAd_WrongText() {
	cAd := testAd
	cAd.Text = wrongEmptyStr

	_, err := s.service.CreateAd(context.Background(), cAd.Title, cAd.Text, cAd.AuthorID)
	assert.ErrorIs(s.T(), ValidationAds.ErrBadText, err)
}

func (s *serviceSuite) Test_AdService_CreateAd_WrongTextMoreLen() {
	cAd := testAd
	cAd.Text = wrongMoreStr

	_, err := s.service.CreateAd(context.Background(), cAd.Title, cAd.Text, cAd.AuthorID)
	assert.ErrorIs(s.T(), ValidationAds.ErrBadText, err)
}

func (s *serviceSuite) Test_AdService_ChangeAdStatus() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	nAd := testAd
	nAd.Published = true
	nAd.UpdateDate = updateDate

	s.adRepo.
		On("EditAdStatus", &cAd, nAd.Published, nAd.UpdateDate).
		Return(&nAd, nil)

	uAd, err := s.service.ChangeAdStatus(context.Background(), testID, cAd.AuthorID, true)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), *uAd, nAd)
}

func (s *serviceSuite) Test_AdService_ChangeAdStatus_WrongID() {
	cAd := testAd

	empty := &entities.Ad{}

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	s.adRepo.
		On("EditAdStatus", &cAd, true, updateDate).
		Return(&empty, nil)

	s.adRepo.
		On("GetAdByID", badID).
		Return(empty, util.ErrNotFound)

	uAd, err := s.service.ChangeAdStatus(context.Background(), badID, cAd.AuthorID, true)

	assert.ErrorIs(s.T(), util.ErrNotFound, err)
	assert.Equal(s.T(), empty, uAd)
}

func (s *serviceSuite) Test_AdService_ChangeAdStatus_WrongAuthorID() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	nAd := testAd
	nAd.Published = true
	nAd.UpdateDate = updateDate

	s.adRepo.
		On("EditAdStatus", cAd, nAd.Published, nAd.UpdateDate).
		Return(&nAd, ValidationAds.ErrBadAuthorID)

	badAuthorID := badID
	uAd, err := s.service.ChangeAdStatus(context.Background(), testID, badAuthorID, true)

	assert.ErrorIs(s.T(), ValidationAds.ErrBadAuthorID, err)
	assert.Equal(s.T(), &cAd, uAd)
}

func (s *serviceSuite) Test_AdService_UpdateAd() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := &entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}

	s.adRepo.
		On("ChangeAdText", cAd.ID, uAd.Title, uAd.Text, updateDate).
		Return(uAd, nil)

	uAd2, err := s.service.UpdateAd(context.Background(), cAd.ID, cAd.AuthorID, uAd.Title, uAd.Text)

	assert.Nil(s.T(), err)
	assert.Equal(s.T(), uAd2, uAd)
}

func (s *serviceSuite) Test_AdService_UpdateAd_WrongID() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}

	empty := &entities.Ad{}

	s.adRepo.
		On("ChangeAdText", badID, uAd.Title, uAd.Text, updateDate).
		Return(empty, util.ErrNotFound)

	s.adRepo.
		On("GetAdByID", badID).
		Return(empty, util.ErrNotFound)

	uAd2, err := s.service.UpdateAd(context.Background(), badID, cAd.AuthorID, uAd.Title, uAd.Text)

	assert.ErrorIs(s.T(), util.ErrNotFound, err)
	assert.Equal(s.T(), uAd2, empty)
}

func (s *serviceSuite) Test_AdService_UpdateAd_WrongAuthorID() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}
	s.adRepo.
		On("ChangeAdText", cAd.ID, uAd.Title, uAd.Text, updateDate).
		Return(&cAd, ValidationAds.ErrBadAuthorID)

	badAuthorID := badID
	uAd2, err := s.service.UpdateAd(context.Background(), testID, badAuthorID, uAd.Title, uAd.Text)

	assert.ErrorIs(s.T(), ValidationAds.ErrBadAuthorID, err)
	assert.Equal(s.T(), *uAd2, cAd)
}

func (s *serviceSuite) Test_AdService_UpdateAd_WrongTitle() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}
	s.adRepo.
		On("ChangeAdText", cAd.ID, uAd.Title, uAd.Text, updateDate).
		Return(&cAd, ValidationAds.ErrBadTitle)

	uAd2, err := s.service.UpdateAd(context.Background(), testID, cAd.AuthorID, wrongEmptyStr, uAd.Text)

	assert.ErrorIs(s.T(), ValidationAds.ErrBadTitle, err)
	assert.Equal(s.T(), *uAd2, cAd)
}

func (s *serviceSuite) Test_AdService_UpdateAd_WrongTitleMoreLen() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}
	s.adRepo.
		On("ChangeAdText", cAd.ID, uAd.Title, uAd.Text, updateDate).
		Return(&cAd, ValidationAds.ErrBadTitle)

	uAd2, err := s.service.UpdateAd(context.Background(), testID, cAd.AuthorID, wrongMoreStr, uAd.Text)

	assert.ErrorIs(s.T(), ValidationAds.ErrBadTitle, err)
	assert.Equal(s.T(), *uAd2, cAd)
}

func (s *serviceSuite) Test_AdService_UpdateAd_WrongText() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}
	s.adRepo.
		On("ChangeAdText", cAd.ID, uAd.Title, uAd.Text, updateDate).
		Return(&cAd, ValidationAds.ErrBadText)

	uAd2, err := s.service.UpdateAd(context.Background(), testID, cAd.AuthorID, uAd.Title, wrongEmptyStr)

	assert.ErrorIs(s.T(), ValidationAds.ErrBadText, err)
	assert.Equal(s.T(), *uAd2, cAd)
}

func (s *serviceSuite) Test_AdService_UpdateAd_WrongTextMoreLen() {
	cAd := testAd

	updateDate, err := s.formatter.ToTime(time.Now().UTC())
	assert.NoError(s.T(), err)

	uAd := entities.Ad{
		ID:         cAd.ID,
		Title:      "newTitle",
		Text:       "newText",
		AuthorID:   cAd.AuthorID,
		CreateDate: cAd.CreateDate,
		UpdateDate: updateDate,
	}
	s.adRepo.
		On("ChangeAdText", cAd.ID, uAd.Title, uAd.Text, updateDate).
		Return(&cAd, ValidationAds.ErrBadText)

	uAd2, err := s.service.UpdateAd(context.Background(), testID, cAd.AuthorID, uAd.Title, wrongMoreStr)

	assert.ErrorIs(s.T(), ValidationAds.ErrBadText, err)
	assert.Equal(s.T(), *uAd2, cAd)
}

func (s *serviceSuite) Test_AdService_RemoveAd() {
	cAd := testAd

	s.adRepo.
		On("DeleteAd", testID).
		Return(nil)

	err := s.service.RemoveAd(context.Background(), cAd.ID, cAd.AuthorID)
	assert.Nil(s.T(), err)
}

func (s *serviceSuite) Test_AdService_RemoveAd_WrongAuthor() {
	cAd := testAd

	s.adRepo.
		On("DeleteAd", testID).
		Return(ValidationAds.ErrBadAuthorID)

	err := s.service.RemoveAd(context.Background(), cAd.ID, 2)
	assert.ErrorIs(s.T(), ValidationAds.ErrBadAuthorID, err)
}

func (s *serviceSuite) Test_AdService_RemoveAd_WrongID() {
	cAd := testAd

	s.adRepo.
		On("DeleteAd", badID).
		Return(util.ErrNotFound)

	empty := &entities.Ad{}

	s.adRepo.
		On("GetAdByID", badID).
		Return(empty, util.ErrNotFound)

	err := s.service.RemoveAd(context.Background(), badID, cAd.AuthorID)
	assert.ErrorIs(s.T(), util.ErrNotFound, err)
}

func (s *serviceSuite) Test_AdService_GetDateTimeFormat() {
	format := s.service.GetDateTimeFormat()
	assert.Equal(s.T(), s.formatter, format)
}

func (s *serviceSuite) Test_AdService_GetAdByID() {
	cAd := testAd

	ad, err := s.service.GetAdByID(context.Background(), cAd.ID)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), cAd, *ad)
}

func Test_AdService_GetAdsByFilter(t *testing.T) {
	AdRepo := new(mocks.AdRepository)
	service := NewAdsService(AdRepo, util.NewDateTimeFormatter(time.DateOnly))

	newAD := testAd
	newAD.Published = true
	expAds := []entities.Ad{newAD}

	AdRepo.
		On("GetAdsByFilters", mock.AnythingOfType("[]func(entities.Ad) bool")).
		Return(expAds, nil)

	ads, err := service.GetAdsByFilter(context.Background(), dFilters)
	assert.Nil(t, err)
	assert.Equal(t, ads, expAds)
}

func (s *serviceSuite) Test_AdService_GetAdsByFilters() {
	cAd := testAd
	cAd.ID = testID
	filters := AdFilters{
		AuthorID:   cAd.AuthorID,
		Published:  cAd.Published,
		CreateDate: cAd.CreateDate,
		Title:      cAd.Title,
	}
	expectedAds := []entities.Ad{cAd}

	s.adRepo.
		On("GetAdsByFilters", mock.AnythingOfType("[]func(entities.Ad) bool")).
		Return(expectedAds, nil)

	ads, err := s.service.GetAdsByFilter(context.Background(), filters)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), ads, expectedAds)
}

func TestGetAdsByFilter(t *testing.T) {
	adRepo := new(mocks.AdRepository)
	service := NewAdsService(adRepo, util.NewDateTimeFormatter(time.DateOnly))

	ad1 := entities.Ad{AuthorID: 1, CreateDate: time.Now(), Title: "Ad 1", Published: true}
	ad2 := entities.Ad{AuthorID: 2, CreateDate: time.Now(), Title: "Ad 2", Published: true}
	expectedAds := []entities.Ad{ad1, ad2}

	filters := AdFilters{}
	adRepo.On("GetAdsByFilters", mock.Anything).Return(expectedAds, nil)
	ads, err := service.GetAdsByFilter(context.Background(), filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedAds, ads)

	filters = AdFilters{AuthorID: 1}
	adRepo.On("GetAdsByFilters", mock.Anything).Return(expectedAds, nil)
	ads, err = service.GetAdsByFilter(context.Background(), filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedAds, ads)

	filters = AdFilters{CreateDate: time.Now()}
	adRepo.On("GetAdsByFilters", mock.Anything).Return(expectedAds, nil)
	ads, err = service.GetAdsByFilter(context.Background(), filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedAds, ads)

	filters = AdFilters{Title: "Ad 1"}
	adRepo.On("GetAdsByFilters", mock.Anything).Return(expectedAds, nil)
	ads, err = service.GetAdsByFilter(context.Background(), filters)
	assert.Nil(t, err)
	assert.Equal(t, expectedAds, ads)
}
