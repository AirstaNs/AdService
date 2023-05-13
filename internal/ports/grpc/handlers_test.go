package grpc

import (
	"context"
	"github.com/AirstaNs/ValidationAds"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"homework10/internal/entities"
	mocks "homework10/internal/mocks/appemocks"
	"homework10/internal/service"
	"homework10/internal/util"
	"strings"
	"testing"
	"time"
)

var (
	emptyUserResp = &UserResponse{}
	emptyAdRem    = &DeleteAdResponse{}
	emptyAdResp   = &AdResponse{}
	emptyUser     = &entities.User{}
	emptyAd       = &entities.Ad{}
	testID        = int64(0)
	badID         = int64(-11111)
	wrongMoreStr  = strings.Repeat("r", 501)
	nPublished    = true
	nTitle        = "newTest"
	nText         = "newTestText"

	tAd = entities.Ad{
		ID:         testID,
		Title:      "Test",
		Text:       "TestText",
		AuthorID:   testID,
		Published:  false,
		CreateDate: time.Now().UTC(),
		UpdateDate: time.Now().UTC(),
	}
	tUser = entities.User{
		ID:       testID,
		Nickname: "Test",
		Email:    "TestEmail",
	}
)

type rpcAppSuite struct {
	suite.Suite
	app  *mocks.App
	serv *GServer
}

func (s *rpcAppSuite) SetupSuite() {
	s.app = new(mocks.App)
	s.serv = &GServer{App: s.app}
}

func (s *rpcAppSuite) SetupTest() {
	s.app = new(mocks.App)
	s.serv = &GServer{App: s.app}

	s.app.
		On("GetUserByID", mock.Anything, tAd.AuthorID).
		Return(&tUser, nil)

	s.app.
		On("CreateAd", mock.Anything, tAd.Title, tAd.Text, tAd.AuthorID).
		Return(&tAd, nil)
}

func TestSuiteRPC(t *testing.T) {
	u := new(rpcAppSuite)
	suite.Run(t, u)

}

func (s *rpcAppSuite) Test_AddAd() {
	background := context.Background()
	ad, err := s.serv.AddAd(background, &CreateAdRequest{
		Title:  tAd.Title,
		Text:   tAd.Text,
		UserId: tAd.AuthorID,
	})
	s.NoError(err)
	s.Equal(AdSuccessResponse(&tAd), ad)

}

func (s *rpcAppSuite) Test_AddAd_NotFoundUser() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()
	app.
		On("GetUserByID", mock.Anything, tAd.AuthorID).
		Return(emptyUser, errNotFound)

	ad, err := s.serv.AddAd(background, &CreateAdRequest{
		Title:  tAd.Title,
		Text:   tAd.Text,
		UserId: tAd.AuthorID,
	})
	s.Error(err)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_AddAd_BadTitle() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()

	app.
		On("GetUserByID", mock.Anything, tAd.AuthorID).
		Return(&tUser, nil)
	app.
		On("CreateAd", mock.Anything, wrongMoreStr, tAd.Text, tAd.AuthorID).
		Return(emptyAd, ValidationAds.ErrBadTitle)

	ad, err := s.serv.AddAd(background, &CreateAdRequest{
		Title:  wrongMoreStr,
		Text:   tAd.Text,
		UserId: tAd.AuthorID,
	})
	s.Error(err, ValidationAds.ErrBadTitle)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_UpdateAdStatus() {
	background := context.Background()

	nAd := tAd
	nAd.Published = nPublished
	nAd.Title = nTitle
	nAd.Text = nText

	s.app.
		On("ChangeAdStatus", mock.Anything, nAd.ID, nAd.AuthorID, nAd.Published).
		Return(&nAd, nil)

	cReq := &ChangeAdStatusRequest{
		AdId:      nAd.ID,
		UserId:    nAd.AuthorID,
		Published: nAd.Published,
	}

	ad, err := s.serv.UpdateAdStatus(background, cReq)
	s.NoError(err)
	s.Equal(AdSuccessResponse(&nAd), ad)
}

func (s *rpcAppSuite) TestUpdateAdStatus_BadAdID() {
	nApp := new(mocks.App)
	s.serv.App = nApp
	background := context.Background()

	cReq := &ChangeAdStatusRequest{
		AdId:      badID,
		UserId:    tAd.AuthorID,
		Published: nPublished,
	}

	nApp.
		On("GetUserByID", mock.Anything, tAd.AuthorID).
		Return(emptyUser, errNotFound)

	ad, err := s.serv.UpdateAdStatus(background, cReq)
	s.Error(err, errNotFound)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_UpdateAdStatus_Forbidden() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()

	fUser := tUser
	fUser.ID = tAd.AuthorID + 1

	cReq := &ChangeAdStatusRequest{
		AdId:      tAd.ID,
		UserId:    fUser.ID,
		Published: nPublished,
	}

	app.
		On("GetUserByID", mock.Anything, fUser.ID).
		Return(&fUser, nil)

	app.
		On("ChangeAdStatus", mock.Anything, cReq.AdId, fUser.ID, cReq.Published).
		Return(&tAd, ValidationAds.ErrBadAuthorID)

	ad, err := s.serv.UpdateAdStatus(background, cReq)
	s.Error(err, errForbidden)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_UpdateAdStatus_BadAdTitle() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()

	nAd := tAd
	nAd.Published = nPublished
	nAd.Title = wrongMoreStr
	nAd.Text = nText

	cReq := &ChangeAdStatusRequest{
		AdId:      nAd.ID,
		UserId:    nAd.AuthorID,
		Published: nAd.Published,
	}

	app.
		On("GetUserByID", mock.Anything, nAd.AuthorID).
		Return(&tUser, nil)

	app.
		On("ChangeAdStatus", mock.Anything, nAd.ID, nAd.AuthorID, nAd.Published).
		Return(emptyAd, ValidationAds.ErrBadTitle)

	ad, err := s.serv.UpdateAdStatus(background, cReq)
	s.Error(err, ValidationAds.ErrBadTitle)
	s.Equal(emptyAdResp, ad)

}

func (s *rpcAppSuite) Test_ModifyAd() {
	background := context.Background()

	nAd := tAd
	nAd.Title = nTitle
	nAd.Text = nText

	s.app.
		On("UpdateAd", mock.Anything, nAd.ID, nAd.AuthorID, nAd.Title, nAd.Text).
		Return(&nAd, nil)

	mReq := &UpdateAdRequest{
		AdId:   nAd.ID,
		UserId: nAd.AuthorID,
		Title:  nAd.Title,
		Text:   nAd.Text,
	}

	ad, err := s.serv.ModifyAd(background, mReq)
	s.NoError(err)
	s.Equal(AdSuccessResponse(&nAd), ad)
}

func (s *rpcAppSuite) Test_ModifyAd_BadAdID() {
	nApp := new(mocks.App)
	s.serv.App = nApp
	background := context.Background()

	mReq := &UpdateAdRequest{
		AdId:   badID,
		UserId: tAd.AuthorID,
		Title:  nTitle,
		Text:   nText,
	}

	nApp.
		On("GetUserByID", mock.Anything, tAd.AuthorID).
		Return(emptyUser, errNotFound)

	ad, err := s.serv.ModifyAd(background, mReq)
	s.Error(err, errNotFound)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_ModifyAd_Forbidden() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()

	fUser := tUser
	fUser.ID = tAd.AuthorID + 1

	mReq := &UpdateAdRequest{
		AdId:   tAd.ID,
		UserId: fUser.ID,
		Title:  nTitle,
		Text:   nText,
	}

	app.
		On("GetUserByID", mock.Anything, fUser.ID).
		Return(&fUser, nil)

	app.
		On("UpdateAd", mock.Anything, mReq.AdId, fUser.ID, mReq.Title, mReq.Text).
		Return(&tAd, ValidationAds.ErrBadAuthorID)

	ad, err := s.serv.ModifyAd(background, mReq)
	s.Error(err, errForbidden)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_ModifyAd_BadAdTitle() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()

	nAd := tAd
	nAd.Title = wrongMoreStr
	nAd.Text = nText

	mReq := &UpdateAdRequest{
		AdId:   nAd.ID,
		UserId: nAd.AuthorID,
		Title:  nAd.Title,
		Text:   nAd.Text,
	}

	app.
		On("GetUserByID", mock.Anything, nAd.AuthorID).
		Return(&tUser, nil)

	app.
		On("UpdateAd", mock.Anything, nAd.ID, nAd.AuthorID, nAd.Title, nAd.Text).
		Return(emptyAd, ValidationAds.ErrBadTitle)

	ad, err := s.serv.ModifyAd(background, mReq)
	s.Error(err, ValidationAds.ErrBadTitle)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_GetAd() {
	background := context.Background()

	s.app.
		On("GetAdByID", mock.Anything, tAd.ID).
		Return(&tAd, nil)

	gReq := &GetADByIDRequest{
		AdId: tAd.ID,
	}

	ad, err := s.serv.GetAd(background, gReq)
	s.NoError(err)
	s.Equal(AdSuccessResponse(&tAd), ad)
}

func (s *rpcAppSuite) Test_GetAd_BadAdID() {
	nApp := new(mocks.App)
	s.serv.App = nApp
	background := context.Background()

	gReq := &GetADByIDRequest{
		AdId: badID,
	}

	nApp.
		On("GetAdByID", mock.Anything, gReq.AdId).
		Return(emptyAd, errNotFound)

	ad, err := s.serv.GetAd(background, gReq)
	s.Error(err, errNotFound)
	s.Equal(emptyAdResp, ad)
}

func (s *rpcAppSuite) Test_RemoveAd() {
	background := context.Background()

	s.app.
		On("RemoveAd", mock.Anything, tAd.ID, tAd.AuthorID).
		Return(nil)

	rReq := &DeleteAdRequest{
		AdId:     tAd.ID,
		AuthorId: tAd.AuthorID,
	}
	response := DeleteAdResponse{AdId: rReq.AdId, UserId: rReq.AuthorId}
	dAdResp, err := s.serv.RemoveAd(background, rReq)
	s.NoError(err)
	s.Equal(&response, dAdResp)
}

func (s *rpcAppSuite) Test_RemoveAd_Forbidden() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()

	fUser := tUser
	fUser.ID = tAd.AuthorID + 1

	rReq := &DeleteAdRequest{
		AdId:     tAd.ID,
		AuthorId: fUser.ID,
	}

	app.
		On("RemoveAd", mock.Anything, rReq.AdId, fUser.ID).
		Return(ValidationAds.ErrBadAuthorID)

	ad, err := s.serv.RemoveAd(background, rReq)
	s.Error(err, errForbidden)
	s.Equal(emptyAdRem, ad)
}

func (s *rpcAppSuite) Test_GetAds() {
	background := context.Background()
	ent := &[]entities.Ad{tAd}
	ListAds := AdListSuccessResponse(ent)
	dFilters := AdFilters{
		OptionalAuthorId:   nil,
		OptionalTitle:      nil,
		OptionalCreateDate: nil,
		OptionalPublished:  nil,
	}
	filters := service.AdFilters{
		AuthorID:   int64(-1),
		CreateDate: time.Time{},
		Title:      "",
		Published:  true,
	}
	s.app.
		On("GetAdsByFilter", mock.Anything, filters).
		Return(*ent, nil)

	s.app.
		On("GetDateTimeFormat").
		Return(util.NewDateTimeFormatter(time.DateOnly), nil)

	ads, err := s.serv.GetAds(background, &dFilters)
	s.NoError(err)
	s.Equal(&ListAds, ads)

}

func (s *rpcAppSuite) Test_GetAds_GoodDate() {
	background := context.Background()
	ent := &[]entities.Ad{tAd}
	ListAds := AdListSuccessResponse(ent)
	dFilters := AdFilters{
		OptionalAuthorId:   nil,
		OptionalTitle:      nil,
		OptionalCreateDate: &timestamppb.Timestamp{Seconds: tAd.CreateDate.Unix()},
		OptionalPublished:  nil,
	}
	s.app.
		On("GetAdsByFilter", mock.Anything, mock.Anything).
		Return(*ent, nil)

	s.app.
		On("GetDateTimeFormat").
		Return(util.NewDateTimeFormatter(time.DateOnly), nil)

	ads, err := s.serv.GetAds(background, &dFilters)
	s.NoError(err)
	s.Equal(&ListAds, ads)

}

func (s *rpcAppSuite) Test_AddUser() {
	background := context.Background()
	uReq := &UserRequest{
		Nickname: tUser.Nickname,
		Email:    tUser.Email,
	}
	s.app.
		On("CreateUser", mock.Anything, tUser.Nickname, tUser.Email).
		Return(&tUser, nil)

	user, err := s.serv.AddUser(background, uReq)
	s.NoError(err)
	s.Equal(UserSuccessResponse(&tUser), user)
}

func (s *rpcAppSuite) Test_GetUser() {
	background := context.Background()
	uReq := &GetUserRequest{
		Id: tUser.ID,
	}
	s.app.
		On("GetUserByID", mock.Anything, tUser.ID).
		Return(&tUser, nil)

	user, err := s.serv.GetUser(background, uReq)
	s.NoError(err)
	s.Equal(UserSuccessResponse(&tUser), user)
}

func (s *rpcAppSuite) Test_GetUser_BadID() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()
	uReq := &GetUserRequest{
		Id: badID,
	}
	app.
		On("GetUserByID", mock.Anything, uReq.Id).
		Return(emptyUser, util.ErrNotFound)

	user, err := s.serv.GetUser(background, uReq)
	s.Error(err, errNotFound)
	s.Equal(emptyUserResp, user)
}

func (s *rpcAppSuite) Test_RemoveUser() {
	background := context.Background()
	uReq := &DeleteUserRequest{
		Id: tUser.ID,
	}
	s.app.
		On("RemoveUser", mock.Anything, tUser.ID).
		Return(nil)

	user, err := s.serv.RemoveUser(background, uReq)
	s.NoError(err)
	s.Equal(&DeleteUserResponse{Id: user.Id}, user)
}

func (s *rpcAppSuite) Test_ModifyUser() {
	background := context.Background()
	uReq := &UserUpdateRequest{
		Id:       tUser.ID,
		Nickname: tUser.Nickname,
		Email:    tUser.Email,
	}
	s.app.
		On("UpdateUser", mock.Anything, tUser.ID, tUser.Nickname, tUser.Email).
		Return(&tUser, nil)

	user, err := s.serv.ModifyUser(background, uReq)
	s.NoError(err)
	s.Equal(UserSuccessResponse(&tUser), user)
}

func (s *rpcAppSuite) Test_ModifyUser_BadID() {
	app := new(mocks.App)
	s.serv.App = app
	background := context.Background()
	uReq := &UserUpdateRequest{
		Id:       badID,
		Nickname: tUser.Nickname,
		Email:    tUser.Email,
	}
	app.
		On("UpdateUser", mock.Anything, uReq.Id, uReq.Nickname, uReq.Email).
		Return(emptyUser, util.ErrNotFound)

	user, err := s.serv.ModifyUser(background, uReq)
	s.Error(err, errNotFound)
	s.Equal(emptyUserResp, user)
}
