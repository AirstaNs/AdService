package gRPC

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"homework10/internal/entities"
	"homework10/internal/ports/grpc"
	"testing"
)

func TestSuiteUsers(t *testing.T) {
	u := new(usersSuite)
	suite.Run(t, u)

}

type usersSuite struct {
	suite.Suite
	client *gRPCtestClient
	users  []entities.User
}

func (s *usersSuite) SetupSuite() {
	s.client = getGRPCTestClient()
	users, err := setupUsers(s.client)
	assert.NoError(s.T(), err)

	s.users = users

}

func (s *usersSuite) TearDownSuite() {
	s.client.Stop()
}

func (s *usersSuite) Test_User_GetByID() {
	server := s.client.Server
	user := s.users[0]

	getUserReq := &grpc.GetUserRequest{Id: user.ID}
	res1, err1 := server.GetUser(context.Background(), getUserReq)
	assert.NoError(s.T(), err1)
	assert.Equal(s.T(), res1.Id, user.ID)
	assert.Equal(s.T(), res1.Nickname, name)
	assert.Equal(s.T(), res1.Email, email)
}

func (s *usersSuite) Test_User_GetByID_WrongID() {
	server := s.client.Server

	getUserReq := &grpc.GetUserRequest{Id: 100}

	_, err1 := server.GetUser(context.Background(), getUserReq)
	assert.Error(s.T(), err1)
}

func (s *usersSuite) Test_User_Create() {
	server := s.client.Server

	userReq := &grpc.UserRequest{Nickname: name, Email: email}
	res, err := server.AddUser(context.Background(), userReq)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res.Nickname, name)
	assert.Equal(s.T(), res.Email, email)
}

func (s *usersSuite) Test_User_Update() {
	server := s.client.Server
	users, err := setupUsers(s.client)
	assert.NoError(s.T(), err)
	user := users[0]

	updateUserReq := &grpc.UserUpdateRequest{Id: user.ID, Nickname: "new name", Email: "new email"}
	res, err := server.ModifyUser(context.Background(), updateUserReq)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), res.Id, user.ID)
}

func (s *usersSuite) Test_User_Delete() {
	server := s.client.Server

	user, err := addUser(s.client, name+email, email+name)
	assert.NoError(s.T(), err)

	deleteUserReq := &grpc.DeleteUserRequest{Id: user.ID}
	_, err = server.RemoveUser(context.Background(), deleteUserReq)
	assert.NoError(s.T(), err)

	deleteUserReq = &grpc.DeleteUserRequest{Id: user.ID}
	_, err = server.RemoveUser(context.Background(), deleteUserReq)
	assert.NoError(s.T(), err)
}
