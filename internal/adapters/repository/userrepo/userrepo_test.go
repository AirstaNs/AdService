package userrepo

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/entities"
	"homework10/internal/util"
	"testing"
)

type repoSuite struct {
	suite.Suite
	repo UserRepository
}

var testUser = entities.User{
	Nickname: "Test",
	Email:    "testuser@example.com",
}

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

func (s *repoSuite) Test_Repo_AddUser() {

	id, err := s.repo.AddUser(testUser)
	assert.NoError(s.T(), err)

	userFromRepo, err := s.repo.GetUserByID(id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *userFromRepo, testUser)
}

func Test_Repo_WrongUID(t *testing.T) {
	rep := &mapRepository{
		rep: make(map[int64]entities.User),
		UID: util.UID{Id: -2}}
	id, err := rep.AddUser(testUser)
	assert.ErrorIs(t, util.ErrGenID, err)
	assert.Equal(t, int64(-1), id)
}

func (s *repoSuite) Test_Repo_GetUserByID_NotFound() {
	_, err := s.repo.GetUserByID(-1)
	assert.ErrorIs(s.T(), err, ErrEmptyUser)
}

func (s *repoSuite) Test_Repo_GetUserByID() {
	id, err := s.repo.AddUser(testUser)
	assert.NoError(s.T(), err)

	newUser := testUser
	newUser.ID = id

	userFromRepo, err := s.repo.GetUserByID(id)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), *userFromRepo, newUser)
}

func (s *repoSuite) Test_Repo_DeleteUser() {
	id, err := s.repo.AddUser(testUser)
	assert.NoError(s.T(), err)

	err = s.repo.DeleteUser(id)
	assert.NoError(s.T(), err)
}

func (s *repoSuite) Test_Repo_DeleteUser_NotFound() {
	err := s.repo.DeleteUser(-1)
	assert.ErrorIs(s.T(), err, ErrEmptyUser)
}

func (s *repoSuite) Test_Repo_UpdateUser() {
	id, err := s.repo.AddUser(testUser)
	assert.NoError(s.T(), err)

	newUser := testUser
	newUser.ID = id
	newUser.Nickname = "NewNickname"
	newUser.Email = "New@mail.ru"

	user, err := s.repo.EditUser(newUser)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), newUser, *user)

}

func (s *repoSuite) Test_Repo_UpdateUser_NotFound() {
	newUser := testUser
	newUser.ID = -1
	newUser.Nickname = "NewNickname"
	newUser.Email = "New@mail.ru"

	_, err := s.repo.EditUser(newUser)
	assert.ErrorIs(s.T(), err, ErrEmptyUser)
}
