package adrepo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/entities"
	"homework10/internal/util"
	"testing"
	"time"
)

type repoSuite struct {
	suite.Suite
	repo AdRepository
}

var dAd = entities.Ad{
	Title:      "Test",
	Text:       "TestText",
	AuthorID:   1,
	Published:  false,
	CreateDate: time.Now().UTC(),
	UpdateDate: time.Now().UTC()}

func TestSuiteAdRepo(t *testing.T) {
	u := new(repoSuite)
	suite.Run(t, u)

}

func (s *repoSuite) SetupSuite() {
	s.repo = New()
}
func (s *repoSuite) TearDownSuite() {
	s.repo = nil
}

func (s *repoSuite) TearDownTest() {
	for key := range s.repo.(*mapRepository).rep {
		delete(s.repo.(*mapRepository).rep, key)
	}
}

func (s *repoSuite) Test_Repo_GetAdByID_NotFound() {
	_, err := s.repo.GetAdByID(-1)
	assert.Error(s.T(), err)
}

func (s *repoSuite) Test_AdRepo_AddAd() {
	id, err := s.repo.AddAd(dAd)
	dId := int64(0)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), dId, id)
}

func Test_AdRepo_GetAdByID_WrongUID(t *testing.T) {
	rep := &mapRepository{
		rep: make(map[int64]entities.Ad),
		UID: util.UID{Id: -2}}
	id, err := rep.AddAd(dAd)
	assert.ErrorIs(t, util.ErrGenID, err)
	assert.Equal(t, int64(-1), id)
}

func (s *repoSuite) Test_AdRepo_EditAdStatus() {
	id, err := s.repo.AddAd(dAd)
	assert.NoError(s.T(), err)

	newAD := dAd
	newAD.ID = id

	var published bool
	updateTime := time.Now().UTC()

	updatedAd, err := s.repo.EditAdStatus(&newAD, published, updateTime)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), published, updatedAd.Published)
	assert.Equal(s.T(), updateTime, updatedAd.UpdateDate)

	adFromRepo, err := s.repo.GetAdByID(id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *updatedAd, *adFromRepo)
}

func (s *repoSuite) Test_AdRepo_ChangeAdText() {
	id, _ := s.repo.AddAd(dAd)
	text := "NewTextUpdate"
	title := "NewTitleUpdate"
	updateTime := time.Now().UTC()
	updatedAd, err := s.repo.ChangeAdText(id, title, text, updateTime)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), title, updatedAd.Title)
	assert.Equal(s.T(), text, updatedAd.Text)
	assert.Equal(s.T(), updateTime, updatedAd.UpdateDate)

	adFromRepo, err := s.repo.GetAdByID(id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *updatedAd, *adFromRepo)
}

func (s *repoSuite) Test_AdRepo_ChangeAdText_WrongAdID() {
	_, _ = s.repo.AddAd(dAd)
	text := "NewTextUpdate"
	title := "NewTitleUpdate"
	updateTime := time.Now().UTC()
	_, err := s.repo.ChangeAdText(-1, title, text, updateTime)
	assert.ErrorIs(s.T(), util.ErrNotFound, err)
}

func (s *repoSuite) Test_AdRepo_GetAdByID() {
	ad, err := s.repo.AddAd(dAd)
	assert.NoError(s.T(), err)
	_, err = s.repo.GetAdByID(ad)
	assert.NoError(s.T(), err)
}

func (s *repoSuite) Test_AdRepo_DeleteAd() {
	id, _ := s.repo.AddAd(dAd)
	err := s.repo.DeleteAd(id)
	assert.NoError(s.T(), err)
}

func (s *repoSuite) Test_AdRepo_GetByFilter_NoFilter() {
	_, err := s.repo.AddAd(dAd)
	assert.NoError(s.T(), err)
	ads, err := s.repo.GetAdsByFilters(nil)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(ads))
}

func (s *repoSuite) Test_AdRepo_GetByFilter_WithAuthorID() {
	_, err := s.repo.AddAd(dAd)
	assert.NoError(s.T(), err)

	newAD := dAd
	newAD.AuthorID = 2
	_, err = s.repo.AddAd(newAD)
	assert.NoError(s.T(), err)

	var adFilters []func(ad entities.Ad) bool
	adFilters = append(adFilters, func(ad entities.Ad) bool {
		return ad.AuthorID == dAd.AuthorID
	})

	ads, err := s.repo.GetAdsByFilters(adFilters)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, len(ads))
}
