package httpgin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AirstaNs/ValidationAds"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/entities"
	mocks "homework10/internal/mocks/appemocks"
	"homework10/internal/service"
	"homework10/internal/util"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

type adResp struct {
	Data adResponse `json:"data"`
}

type adListResp struct {
	Data []adResponse `json:"data"`
}

var (
	emptyUser    = &entities.User{}
	emptyAd      = &entities.Ad{}
	testID       = int64(0)
	badID        = int64(-11111)
	wrongMoreStr = strings.Repeat("r", 501)
	nPublished   = true
	nTitle       = "newTest"
	nText        = "newTestText"

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

	dFilters = service.AdFilters{
		AuthorID:   int64(-1),
		CreateDate: time.Time{},
		Title:      "",
		Published:  true,
	}
)

type httpAppSuite struct {
	suite.Suite
	app      *mocks.App
	ctx      *gin.Context
	recorder *httptest.ResponseRecorder
}

func TestSuiteHttpApp(t *testing.T) {
	u := new(httpAppSuite)
	suite.Run(t, u)

}

func (s *httpAppSuite) SetupSuite() {
	mApp := new(mocks.App)
	mApp.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), testID).
		Return(&tUser, nil)
	mApp.
		On("CreateAd", mock.AnythingOfType("*gin.Context"), tAd.Title, tAd.Text, tAd.ID).
		Return(&tAd, nil)

	s.app = mApp

}

func (s *httpAppSuite) SetupTest() {
	w := httptest.NewRecorder()
	s.ctx = GetTestGinContext(w)
	s.recorder = w

}

func (s *httpAppSuite) Test_CreateAd() {

	body := map[string]any{
		"user_id": tAd.AuthorID,
		"title":   tAd.Title,
		"text":    tAd.Text,
	}

	MockJsonPost(s.ctx, body)
	createAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusCreated, s.recorder.Code)

	response, err := getResponse(s.recorder)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), tAd, unmarshalResponse(response))
}

func (s *httpAppSuite) Test_CreateAd_InvalidBody() {

	body := map[string]any{
		"user_id": tAd.Text,
	}

	MockJsonPost(s.ctx, body)
	createAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_CreateAd_InvalidUserID() {

	body := map[string]any{
		"user_id": badID,
		"title":   tAd.Title,
		"text":    tAd.Text,
	}
	s.app.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), badID).
		Return(emptyUser, userrepo.ErrEmptyUser)

	MockJsonPost(s.ctx, body)
	createAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusNotFound, s.recorder.Code)
}

func (s *httpAppSuite) Test_CreateAd_CreateAd_BadTitle() {

	nAd := tAd
	nAd.Title = wrongMoreStr

	body := map[string]any{
		"user_id": nAd.AuthorID,
		"title":   nAd.Title,
		"text":    nAd.Text,
	}
	s.app.
		On("CreateAd", mock.AnythingOfType("*gin.Context"), nAd.Title, tAd.Text, tAd.ID).
		Return(&nAd, ValidationAds.ErrBadTitle)

	MockJsonPost(s.ctx, body)
	createAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_CreateAd_CreateAd_BadText() {
	nAd := tAd
	nAd.Text = wrongMoreStr

	body := map[string]any{
		"user_id": nAd.AuthorID,
		"title":   nAd.Title,
		"text":    nAd.Text,
	}
	s.app.
		On("CreateAd", mock.AnythingOfType("*gin.Context"), nAd.Title, nAd.Text, nAd.AuthorID).
		Return(&nAd, ValidationAds.ErrBadText)

	MockJsonPost(s.ctx, body)
	createAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_ChangeAdStatus() {
	body := map[string]any{
		"user_id":   tUser.ID,
		"published": nPublished,
	}
	nAd := tAd
	nAd.Published = nPublished
	s.app.
		On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), tAd.ID, tAd.AuthorID, nPublished).
		Return(&nAd, nil)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	changeAdStatus(s.app)(s.ctx)

	response, err := getResponse(s.recorder)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), nAd, unmarshalResponse(response))
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_ChangeAdStatus_InvalidBody() {
	body := map[string]any{
		"user_id": "кккк",
	}
	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	changeAdStatus(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_ChangeAdStatus_InvalidAdID() {
	body := map[string]any{
		"user_id":   tUser.ID,
		"published": nPublished,
	}
	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: wrongMoreStr}})
	changeAdStatus(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_ChangeAdStatus_InvalidUserID() {
	body := map[string]any{
		"user_id":   badID,
		"published": nPublished,
	}
	s.app.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), badID).
		Return(emptyUser, userrepo.ErrEmptyUser)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	changeAdStatus(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusForbidden, s.recorder.Code)
}

func (s *httpAppSuite) Test_ChangeAdStatus_InvalidAdUserIDForbidden() {
	mApp := new(mocks.App)
	body := map[string]any{
		"user_id":   tUser.ID,
		"published": nPublished,
	}

	mApp.
		On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), tAd.ID, tAd.AuthorID, nPublished).
		Return(emptyAd, ValidationAds.ErrBadAuthorID)

	mApp.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), testID).
		Return(&tUser, nil)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	changeAdStatus(mApp)(s.ctx)
	assert.EqualValues(s.T(), http.StatusForbidden, s.recorder.Code)
}

func (s *httpAppSuite) Test_ChangeAdStatus_InvalidTitle() {
	mApp := new(mocks.App)
	body := map[string]any{
		"user_id":   tUser.ID,
		"published": nPublished,
	}
	nAd := tAd
	nAd.Title = wrongMoreStr
	mApp.
		On("ChangeAdStatus", mock.AnythingOfType("*gin.Context"), tAd.ID, tAd.AuthorID, nPublished).
		Return(&nAd, ValidationAds.ErrBadTitle)

	mApp.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), testID).
		Return(&tUser, nil)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	changeAdStatus(mApp)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_UpdateAd() {
	nAd := tAd
	nAd.Title = nTitle
	nAd.Text = nText

	body := map[string]any{
		"user_id": tUser.ID,
		"title":   nAd.Title,
		"text":    nAd.Text,
	}
	s.app.
		On("GetAdByID", mock.AnythingOfType("*gin.Context"), tAd.ID).
		Return(&tAd, nil)

	s.app.
		On("UpdateAd", mock.AnythingOfType("*gin.Context"), tAd.ID, tAd.AuthorID, nTitle, nText).
		Return(&nAd, nil)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	updateAd(s.app)(s.ctx)

	response, err := getResponse(s.recorder)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), nAd, unmarshalResponse(response))
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)
}

func (s *httpAppSuite) Test_UpdateAdInvalidBody() {
	body := map[string]any{
		"user_id": "кккк",
	}
	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	updateAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_UpdateAd_InvalidAdID() {
	body := map[string]any{}

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: wrongMoreStr}})
	updateAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_UpdateAd_InvalidAdID_NotFound() {
	mApp := new(mocks.App)
	body := map[string]any{}

	mApp.
		On("GetAdByID", mock.AnythingOfType("*gin.Context"), badID).
		Return(emptyAd, util.ErrNotFound)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(badID, 10)}})
	updateAd(mApp)(s.ctx)
	assert.EqualValues(s.T(), http.StatusNotFound, s.recorder.Code)
}

func (s *httpAppSuite) Test_UpdateAd_InvalidUserID() {
	mApp := new(mocks.App)
	body := map[string]any{
		"user_id": badID,
		"title":   nTitle,
		"text":    nText,
	}
	mApp.
		On("GetAdByID", mock.AnythingOfType("*gin.Context"), tAd.ID).
		Return(&tAd, nil)

	mApp.
		On("UpdateAd", mock.AnythingOfType("*gin.Context"), tAd.ID, badID, nTitle, nText).
		Return(emptyAd, ValidationAds.ErrBadAuthorID)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	updateAd(mApp)(s.ctx)
	assert.EqualValues(s.T(), http.StatusForbidden, s.recorder.Code)
}

func (s *httpAppSuite) Test_UpdateAd_InvalidTitle() {
	mApp := new(mocks.App)
	body := map[string]any{
		"user_id": tUser.ID,
		"title":   wrongMoreStr,
		"text":    nText,
	}
	mApp.
		On("GetAdByID", mock.AnythingOfType("*gin.Context"), tAd.ID).
		Return(&tAd, nil)

	mApp.
		On("UpdateAd", mock.AnythingOfType("*gin.Context"), tAd.ID, tAd.AuthorID, wrongMoreStr, nText).
		Return(emptyAd, ValidationAds.ErrBadTitle)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}})
	updateAd(mApp)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_getAdByID() {
	s.app.
		On("GetAdByID", mock.AnythingOfType("*gin.Context"), tAd.ID).
		Return(&tAd, nil)

	MockJsonGet(s.ctx, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}}, nil)
	getAdByID(s.app)(s.ctx)

	response, err := getResponse(s.recorder)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), tAd, unmarshalResponse(response))
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_getAdByID_InvalidAdID() {
	MockJsonGet(s.ctx, gin.Params{{Key: "ad_id", Value: wrongMoreStr}}, nil)
	getAdByID(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_getAdByID_InvalidAdID_NotFound() {
	s.app.
		On("GetAdByID", mock.AnythingOfType("*gin.Context"), badID).
		Return(emptyAd, util.ErrNotFound)

	MockJsonGet(s.ctx, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(badID, 10)}}, nil)
	getAdByID(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusNotFound, s.recorder.Code)
}

func (s *httpAppSuite) Test_DeleteAd() {
	s.app.
		On("RemoveAd", mock.AnythingOfType("*gin.Context"), tAd.ID, tUser.ID).
		Return(nil)

	MockJsonDelete(s.ctx, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}}, url.Values{
		"user_id": []string{strconv.FormatInt(tUser.ID, 10)},
	})
	deleteAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_DeleteAd_InvalidAdID() {
	MockJsonDelete(s.ctx, gin.Params{{Key: "ad_id", Value: wrongMoreStr}}, url.Values{
		"user_id": []string{strconv.FormatInt(tUser.ID, 10)},
	})
	deleteAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_DeleteAd_InvalidUserID() {
	s.app.
		On("RemoveAd", mock.AnythingOfType("*gin.Context"), tAd.ID, badID).
		Return(ValidationAds.ErrBadAuthorID)

	MockJsonDelete(s.ctx, gin.Params{{Key: "ad_id", Value: strconv.FormatInt(tAd.ID, 10)}}, url.Values{
		"user_id": []string{strconv.FormatInt(badID, 10)},
	})
	deleteAd(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusForbidden, s.recorder.Code)
}

func (s *httpAppSuite) Test_getAdsByFilter() {
	newAD := tAd
	newAD.Published = true
	expAds := []entities.Ad{newAD}

	s.app.
		On("GetAdsByFilter", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("service.AdFilters")).
		Return(expAds, nil)

	values := url.Values{
		"user_id":   {strconv.FormatInt(dFilters.AuthorID, 10)},
		"published": {strconv.FormatBool(dFilters.Published)},
	}
	fmt.Println(values.Encode())
	MockJsonGet(s.ctx, gin.Params{}, values)
	getAdsByFilter(s.app)(s.ctx)

	response, err := getResponseList(s.recorder)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), expAds, unmarshalResponseList(response))
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_getAdsByFilter_InvalidUserID() {
	MockJsonGet(s.ctx, gin.Params{}, url.Values{
		"user_id":   {wrongMoreStr},
		"published": {strconv.FormatBool(dFilters.Published)},
	})
	getAdsByFilter(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_getAdsByFilter_InvalidFilter() {
	mApp := new(mocks.App)
	mApp.
		On("GetAdsByFilter", mock.AnythingOfType("*gin.Context"), mock.AnythingOfType("service.AdFilters")).
		Return(nil, ValidationAds.ErrBadText)

	MockJsonGet(s.ctx, gin.Params{}, url.Values{
		"user_id":   {strconv.FormatInt(dFilters.AuthorID, 10)},
		"published": {strconv.FormatBool(dFilters.Published)},
		"text":      {wrongMoreStr},
	})
	getAdsByFilter(mApp)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_createUser() {
	s.app.
		On("CreateUser", mock.AnythingOfType("*gin.Context"), tUser.Nickname, tUser.Email).
		Return(&tUser, nil)

	MockJsonPost(s.ctx, tUser)
	createUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusCreated, s.recorder.Code)

}

func (s *httpAppSuite) Test_createUser_InvalidUser() {
	body := map[string]any{
		"nickname": time.Now().UTC().String(),
	}
	MockJsonPost(s.ctx, len(body))
	createUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_getUserByID() {
	s.app.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), tUser.ID).
		Return(&tUser, nil)

	MockJsonGet(s.ctx, gin.Params{{Key: "user_id", Value: strconv.FormatInt(tUser.ID, 10)}}, url.Values{})
	getUserByID(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_getUserByID_InvalidUserID() {
	MockJsonGet(s.ctx, gin.Params{{Key: "user_id", Value: wrongMoreStr}}, url.Values{})
	getUserByID(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_getUserByID_NotFound() {
	s.app.
		On("GetUserByID", mock.AnythingOfType("*gin.Context"), badID).
		Return(emptyUser, userrepo.ErrEmptyUser)

	MockJsonGet(s.ctx, gin.Params{{Key: "user_id", Value: strconv.FormatInt(badID, 10)}}, url.Values{})
	getUserByID(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusNotFound, s.recorder.Code)
}

func (s *httpAppSuite) Test_updateUser() {
	nUser := tUser
	nUser.Nickname = "newNickname"
	nUser.Email = "newEmail"
	body := map[string]any{
		"nickname": nUser.Nickname,
		"email":    nUser.Email,
	}
	s.app.
		On("UpdateUser", mock.AnythingOfType("*gin.Context"), nUser.ID, nUser.Nickname, nUser.Email).
		Return(&tUser, nil)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "user_id", Value: strconv.FormatInt(tUser.ID, 10)}})
	updateUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_updateUser_InvalidUserID() {
	MockJsonPut(s.ctx, 4, gin.Params{{Key: "user_id", Value: wrongMoreStr}})
	updateUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_updateUser_InvalidID() {
	body := map[string]any{
		"nickname": tUser.Nickname,
		"email":    tUser.Email,
	}
	s.app.
		On("UpdateUser", mock.AnythingOfType("*gin.Context"), badID, tUser.Nickname, tUser.Email).
		Return(emptyUser, userrepo.ErrEmptyUser)

	MockJsonPut(s.ctx, body, gin.Params{{Key: "user_id", Value: strconv.FormatInt(badID, 10)}})
	updateUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func (s *httpAppSuite) Test_deleteUser() {
	s.app.
		On("RemoveUser", mock.AnythingOfType("*gin.Context"), tUser.ID).
		Return(nil)

	MockJsonDelete(s.ctx, gin.Params{{Key: "user_id", Value: strconv.FormatInt(tUser.ID, 10)}}, nil)
	deleteUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusOK, s.recorder.Code)

}

func (s *httpAppSuite) Test_deleteUser_InvalidUserID() {
	MockJsonDelete(s.ctx, gin.Params{{Key: "user_id", Value: wrongMoreStr}}, nil)
	deleteUser(s.app)(s.ctx)
	assert.EqualValues(s.T(), http.StatusBadRequest, s.recorder.Code)
}

func MockJsonGet(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "GET"
	c.Request.Header.Set("Content-Type", "application/json")

	// set path params
	c.Params = params
	c.Request.URL.RawQuery = u.Encode()
}

func MockJsonDelete(c *gin.Context, params gin.Params, u url.Values) {
	c.Request.Method = "DELETE"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = params
	c.Request.URL.RawQuery = u.Encode()
}

func MockJsonPost(c *gin.Context, content any) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	response, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(response))
}

func MockJsonPut(c *gin.Context, content any, params gin.Params) {
	c.Request.Method = "PUT"
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = params

	response, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(response))
}

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func getResponse(w *httptest.ResponseRecorder) (adResp, error) {
	emptyResp := adResp{}
	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		return emptyResp, err
	}

	var response adResp
	if err = json.Unmarshal(respBody, &response); err != nil {
		return emptyResp, err
	}
	return response, nil
}

func getResponseList(w *httptest.ResponseRecorder) (adListResp, error) {
	emptyResp := adListResp{}
	respBody, err := io.ReadAll(w.Body)
	if err != nil {
		return emptyResp, err
	}

	var response adListResp
	if err = json.Unmarshal(respBody, &response); err != nil {
		return emptyResp, err
	}
	return response, nil
}

func unmarshalResponse(response adResp) entities.Ad {
	return entities.Ad{
		ID:         response.Data.ID,
		Title:      response.Data.Title,
		Text:       response.Data.Text,
		AuthorID:   response.Data.AuthorID,
		Published:  response.Data.Published,
		CreateDate: response.Data.CreateDate,
		UpdateDate: response.Data.UpdateDate,
	}
}

func unmarshalResponseList(response adListResp) []entities.Ad {
	var ads []entities.Ad
	for _, ad := range response.Data {
		ads = append(ads, entities.Ad{
			ID:         ad.ID,
			Title:      ad.Title,
			Text:       ad.Text,
			AuthorID:   ad.AuthorID,
			Published:  ad.Published,
			CreateDate: ad.CreateDate,
			UpdateDate: ad.UpdateDate,
		})
	}
	return ads
}
