package gRPC

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"homework10/internal/entities"
	"homework10/internal/ports/grpc"
	"math"
	"testing"
)

type adsSuite struct {
	suite.Suite
	client *gRPCtestClient
	users  []entities.User
	ads    []entities.Ad
}

func (s *adsSuite) SetupSuite() {
	s.client = getGRPCTestClient()
	users, err := setupUsers(s.client)
	assert.NoError(s.T(), err)

	s.users = users

	ads, err := setupAds(s.client, s.users[0], s.users[1])
	assert.NoError(s.T(), err)
	s.ads = ads
}

func (s *adsSuite) TearDownSuite() {
	s.client.Stop()
}

func (s *adsSuite) Test_Ads_GetByID() {
	server := s.client.Server
	ad := s.ads[0]

	getAdReq := &grpc.GetADByIDRequest{AdId: ad.ID}
	res2, err2 := server.GetAd(context.Background(), getAdReq)
	assert.NoError(s.T(), err2)
	assert.Equal(s.T(), res2.Id, ad.ID)
	assert.Equal(s.T(), res2.Title, ad.Title)
	assert.Equal(s.T(), res2.Text, ad.Text)
}

func (s *adsSuite) Test_Ads_GetByID_NoExistID() {
	server := s.client.Server

	getAdReq := &grpc.GetADByIDRequest{AdId: 100}
	_, err2 := server.GetAd(context.Background(), getAdReq)
	assert.ErrorIs(s.T(), err2, errNotFound)
}

func (s *adsSuite) Test_Ads_GetByFilter_NoFilter() {
	server := s.client.Server

	ad := s.ads[0]
	ad1 := s.ads[1]

	sChange := grpc.ChangeAdStatusRequest{AdId: ad.ID, UserId: ad.AuthorID, Published: true}
	_, err := server.UpdateAdStatus(context.Background(), &sChange)
	assert.NoError(s.T(), err)

	sChange2 := grpc.ChangeAdStatusRequest{AdId: ad1.ID, UserId: ad1.AuthorID, Published: true}
	_, err = server.UpdateAdStatus(context.Background(), &sChange2)
	assert.NoError(s.T(), err)

	filters := grpc.AdFilters{}
	listAds, err := server.GetAds(context.Background(), &filters)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(listAds.List), 2)

}

func (s *adsSuite) Test_Ads_GetByFilter_WithTitle() {
	server := s.client.Server

	ad := s.ads[0]
	ad1 := s.ads[1]

	sChange := grpc.ChangeAdStatusRequest{AdId: ad.ID, UserId: ad.AuthorID, Published: true}
	_, err := server.UpdateAdStatus(context.Background(), &sChange)
	assert.NoError(s.T(), err)

	sChange2 := grpc.ChangeAdStatusRequest{AdId: ad1.ID, UserId: ad1.AuthorID, Published: false}
	_, err = server.UpdateAdStatus(context.Background(), &sChange2)
	assert.NoError(s.T(), err)
	titleFilter := wrapperspb.String(title)
	filters := grpc.AdFilters{OptionalTitle: titleFilter}

	listAds, err := server.GetAds(context.Background(), &filters)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(listAds.List), 1)
	assert.Equal(s.T(), listAds.List[0].Title, ad.Title)

}

func (s *adsSuite) Test_Ads_GetByFilter_WithAuthorID() {
	server := s.client.Server

	ad := s.ads[0]
	ad1 := s.ads[1]
	user := s.users[0]

	sChange := grpc.ChangeAdStatusRequest{AdId: ad.ID, UserId: ad.AuthorID, Published: true}
	_, err := server.UpdateAdStatus(context.Background(), &sChange)
	assert.NoError(s.T(), err)

	sChange2 := grpc.ChangeAdStatusRequest{AdId: ad1.ID, UserId: ad1.AuthorID, Published: false}
	_, err = server.UpdateAdStatus(context.Background(), &sChange2)
	assert.NoError(s.T(), err)

	AuthorIdFilters := wrapperspb.Int64(user.ID)
	filters := grpc.AdFilters{OptionalAuthorId: AuthorIdFilters}
	listAds, err := server.GetAds(context.Background(), &filters)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(listAds.List), 1)
	assert.Equal(s.T(), listAds.List[0].AuthorId, user.ID)

}

func (s *adsSuite) Test_Ads_GetByFilter_WithAllFilers() {
	server := s.client.Server
	ad := s.ads[0]
	ad1 := s.ads[1]
	user := s.users[0]

	sChange := grpc.ChangeAdStatusRequest{AdId: ad.ID, UserId: ad.AuthorID, Published: true}
	updateAd := setupUpdateAd(s.client, &sChange)
	ad = *updateAd

	sChange1 := grpc.ChangeAdStatusRequest{AdId: ad1.ID, UserId: ad1.AuthorID, Published: false}
	updateAd1 := setupUpdateAd(s.client, &sChange1)
	ad1 = *updateAd1

	titleFilter := wrapperspb.String(ad.Title)
	authIdFilter := wrapperspb.Int64(user.ID)
	cDateFilter := timestamppb.New(ad.CreateDate)
	publishedFilter := wrapperspb.Bool(ad.Published)

	filters := grpc.AdFilters{
		OptionalTitle:      titleFilter,
		OptionalAuthorId:   authIdFilter,
		OptionalCreateDate: cDateFilter,
		OptionalPublished:  publishedFilter,
	}

	listAds, err := server.GetAds(context.Background(), &filters)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(listAds.List), 1)
	assert.Equal(s.T(), listAds.List[0].Title, ad.Title)
	assert.Equal(s.T(), listAds.List[0].AuthorId, user.ID)
	assert.Equal(s.T(), listAds.List[0].CreateDate.AsTime().UTC(), ad.CreateDate)
	assert.Equal(s.T(), listAds.List[0].Published, ad.Published)
	assert.Equal(s.T(), listAds.List[0].Id, ad.ID)
	assert.Equal(s.T(), listAds.List[0].UpdateDate.AsTime().UTC(), ad.UpdateDate)

}

func (s *adsSuite) Test_Ads_GetByFilter_WithWrongFilters() {
	server := s.client.Server
	ad := s.ads[0]
	ad1 := s.ads[1]
	user := s.users[0]

	sChange := grpc.ChangeAdStatusRequest{AdId: ad.ID, UserId: ad.AuthorID, Published: true}
	updateAd := setupUpdateAd(s.client, &sChange)
	ad = *updateAd

	sChange1 := grpc.ChangeAdStatusRequest{AdId: ad1.ID, UserId: ad1.AuthorID, Published: false}
	updateAd1 := setupUpdateAd(s.client, &sChange1)
	ad1 = *updateAd1

	titleFilter := wrapperspb.String("wrong title")
	authIdFilter := wrapperspb.Int64(user.ID)
	cDateFilter := timestamppb.New(ad.CreateDate)
	publishedFilter := wrapperspb.Bool(ad.Published)

	filters := grpc.AdFilters{
		OptionalTitle:      titleFilter,
		OptionalAuthorId:   authIdFilter,
		OptionalCreateDate: cDateFilter,
		OptionalPublished:  publishedFilter,
	}

	listAds, err := server.GetAds(context.Background(), &filters)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), len(listAds.List), 0)

}

func (s *adsSuite) Test_Ads_Create() {
	newTitle := title + title
	newText := text + text

	ad, err := addAd(s.client, newTitle, newText, s.users[0].ID)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), ad.Title, newTitle)
	assert.Equal(s.T(), ad.Text, newText)
	_, err = s.client.Server.RemoveAd(context.Background(), &grpc.DeleteAdRequest{AdId: ad.ID, AuthorId: ad.AuthorID})
	assert.NoError(s.T(), err)
}

func (s *adsSuite) Test_Ads_Create_WithWrongAuthorID() {
	newTitle := title + title
	newText := text + text

	_, err2 := addAd(s.client, newTitle, newText, math.MaxInt)
	assert.ErrorIs(s.T(), err2, errNotFound)
}

func (s *adsSuite) Test_Ads_Update() {
	server := s.client.Server
	ad := s.ads[0]
	newTitle := title + title
	newText := text + text

	updateAdReq := &grpc.UpdateAdRequest{
		AdId:   ad.ID,
		UserId: ad.AuthorID,
		Title:  newTitle,
		Text:   newText,
	}

	updateAd, err := server.ModifyAd(context.Background(), updateAdReq)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), updateAd.Title, newTitle)
	assert.Equal(s.T(), updateAd.Text, newText)

	_, err = s.client.Server.RemoveAd(context.Background(), &grpc.DeleteAdRequest{AdId: ad.ID, AuthorId: ad.AuthorID})
	assert.NoError(s.T(), err)
}

func (s *adsSuite) Test_Ads_UpdateStatus() {
	user, err := addUser(s.client, "test", "test@mail.ru")
	assert.NoError(s.T(), err)

	ad, err := addAd(s.client, title+title, text+text, user.ID)
	assert.NoError(s.T(), err)

	sChange := grpc.ChangeAdStatusRequest{AdId: ad.ID, UserId: ad.AuthorID, Published: true}
	updateAd := setupUpdateAd(s.client, &sChange)
	assert.Equal(s.T(), updateAd.Published, true)

	_, err = s.client.Server.RemoveAd(context.Background(), &grpc.DeleteAdRequest{AdId: ad.ID, AuthorId: user.ID})
	assert.NoError(s.T(), err)

}

func (s *adsSuite) Test_Ads_Delete_Forbidden() {
	server := s.client.Server
	ad := s.ads[0]
	deleteAdReq := &grpc.DeleteAdRequest{AdId: ad.ID, AuthorId: math.MaxInt}
	_, err := server.RemoveAd(context.Background(), deleteAdReq)
	assert.ErrorIs(s.T(), err, errForbidden)
}

func (s *adsSuite) Test_Ads_Delete() {
	server := s.client.Server

	ad, err := addAd(s.client, title+title, text+text, s.users[0].ID)
	assert.NoError(s.T(), err)

	deleteAdReq := &grpc.DeleteAdRequest{AdId: ad.ID, AuthorId: ad.AuthorID}
	_, err = server.RemoveAd(context.Background(), deleteAdReq)
	assert.NoError(s.T(), err)

	deleteAdReq = &grpc.DeleteAdRequest{AdId: ad.ID, AuthorId: ad.AuthorID}
	_, err = server.RemoveAd(context.Background(), deleteAdReq)
	assert.NoError(s.T(), err)

	_, err = server.GetAd(context.Background(), &grpc.GetADByIDRequest{AdId: ad.ID})
	assert.ErrorIs(s.T(), err, errNotFound)
}

func TestSuiteAds(t *testing.T) {
	a := new(adsSuite)
	suite.Run(t, a)
}
