package service

import (
	"context"
	"github.com/stretchr/testify/suite"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/entities"
	mocks "homework10/internal/mocks/repomocks"
	"testing"
)

type serviceSuiteUsers struct {
	suite.Suite
	service UserService
	uRepo   *mocks.UserRepository
}

var (
	badUserID  = int64(-4124)
	testUserID = int64(0)
	tUser      = entities.User{
		Nickname: "test",
		Email:    "test@mail.ru",
	}
	emptyUser = &entities.User{}
)

func TestSuiteUserService(t *testing.T) {
	u := new(serviceSuiteUsers)
	suite.Run(t, u)

}

func (s *serviceSuiteUsers) SetupSuite() {
	userRepo := new(mocks.UserRepository)
	s.service = NewUserService(userRepo)
	s.uRepo = userRepo

	tUser.ID = testUserID
	s.uRepo.
		On("AddUser", tUser).
		Return(testUserID, nil)

	s.uRepo.
		On("GetUserByID", testUserID).
		Return(&tUser, nil)
}

func (s *serviceSuiteUsers) TearDownSuiteUser() {
	s.service = nil
}

func (s *serviceSuiteUsers) TestCreateUser() {
	aUser := tUser

	user, err := s.service.CreateUser(context.Background(), tUser.Nickname, tUser.Email)
	s.Nil(err)
	s.Equal(&aUser, user)
}

func (s *serviceSuiteUsers) TestGetUserByID() {
	aUser := tUser

	user, err := s.service.GetUserByID(context.Background(), testUserID)
	s.Nil(err)
	s.Equal(&aUser, user)
}

func (s *serviceSuiteUsers) TestUpdateUser() {
	uUser := tUser
	uUser.Nickname = "testNew"
	uUser.Email = "testNew@mail.ru"

	s.uRepo.
		On("EditUser", uUser).
		Return(&uUser, nil)

	user, err := s.service.UpdateUser(context.Background(), testUserID, uUser.Nickname, uUser.Email)
	s.Nil(err)
	s.Equal(&uUser, user)
}

func (s *serviceSuiteUsers) TestUpdateUserWrongID() {
	uUser := tUser
	uUser.Nickname = "testNew"
	uUser.Email = "testNew@mail.ru"

	s.uRepo.
		On("EditUser", tUser).
		Return(&uUser, nil)

	s.uRepo.
		On("GetUserByID", badUserID).
		Return(emptyUser, userrepo.ErrEmptyUser)

	user, err := s.service.UpdateUser(context.Background(), badUserID, uUser.Nickname, uUser.Email)
	s.ErrorIs(userrepo.ErrEmptyUser, err)
	s.Equal(emptyUser, user)
}

func (s *serviceSuiteUsers) TestRemoveUser() {
	s.uRepo.
		On("DeleteUser", testUserID).
		Return(nil)

	err := s.service.RemoveUser(context.Background(), testUserID)
	s.Nil(err)
}

func (s *serviceSuiteUsers) TestRemoveUserWrongID() {
	s.uRepo.
		On("DeleteUser", badUserID).
		Return(userrepo.ErrEmptyUser)

	s.uRepo.
		On("GetUserByID", badUserID).
		Return(emptyUser, userrepo.ErrEmptyUser)

	err := s.service.RemoveUser(context.Background(), badUserID)
	s.ErrorIs(userrepo.ErrEmptyUser, err)
}

func BenchmarkUsersService_CreateUser(b *testing.B) {
	uRepo := new(mocks.UserRepository)
	service := NewUserService(uRepo)

	nUser := tUser
	tUser.ID = testUserID

	uRepo.
		On("AddUser", nUser).
		Return(testUserID, nil)

	for i := 0; i < b.N; i++ {
		_, _ = service.CreateUser(context.Background(), nUser.Nickname, nUser.Email)
	}
}
