package http

import (
	"github.com/stretchr/testify/assert"
	"homework10/internal/util"
	"strconv"
	"testing"
	"time"
)

func Test_Ads_GetByID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	response, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.getAdByID(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, userID)
	assert.False(t, response.Data.Published)
}

func Test_Ads_GetByID_NoExistID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	_, err = client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.getAdByID(100)
	assert.ErrorIs(t, err, ErrorNotFound)
}

func Test_Ads_GetByFilter_NoFilter(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	response, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	response1, err := client.createAd(user.Data.ID, "world", "hello")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user.Data.ID, response1.Data.ID, true)
	assert.NoError(t, err)

	responseList, err := client.listAds()
	assert.NoError(t, err)
	assert.Equal(t, len(responseList.Data), 2)
}

func Test_Ads_GetByFilter_WithAuthorID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	user1, err := client.createUser("qwertys1", "@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "world", "hello")
	assert.NoError(t, err)

	authorFilter := queryParam{"user_id": strconv.Itoa(int(user.Data.ID))}
	responseList, err := client.listAdsFilters(authorFilter)
	assert.NoError(t, err)
	assert.Equal(t, len(responseList.Data), 1)
	assert.Equal(t, responseList.Data[0].AuthorID, user.Data.ID)
}

func Test_Ads_GetByFilter_WithTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	user1, err := client.createUser("qwertys1", "@mail.ru")
	assert.NoError(t, err)

	respone, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "world", "hello")
	assert.NoError(t, err)

	titleFilter := queryParam{"title": respone.Data.Title}
	responseList, err := client.listAdsFilters(titleFilter)
	assert.NoError(t, err)
	assert.Equal(t, len(responseList.Data), 1)
	assert.Equal(t, responseList.Data[0].Title, respone.Data.Title)
}

func Test_Ads_GetByFilter_WithAllFilers(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	user1, err := client.createUser("qwertys1", "@mail.ru")
	assert.NoError(t, err)

	respone, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "world", "hello")
	assert.NoError(t, err)

	s := respone.Data.CreateDate

	formatter := util.NewDateTimeFormatter(time.RFC3339)
	fromString := formatter.ToString(s)

	allFilters := queryParam{
		"published":   strconv.FormatBool(respone.Data.Published),
		"user_id":     strconv.Itoa(int(user.Data.ID)),
		"title":       respone.Data.Title,
		"create_Date": fromString,
	}

	responseList, err := client.listAdsFilters(allFilters)
	assert.NoError(t, err)
	assert.Equal(t, len(responseList.Data), 1)
	assert.Equal(t, responseList.Data[0].Title, respone.Data.Title)
	assert.Equal(t, responseList.Data[0].AuthorID, user.Data.ID)
	assert.Equal(t, responseList.Data[0].Published, respone.Data.Published)
	equalDate := responseList.Data[0].CreateDate.Equal(respone.Data.CreateDate)
	assert.True(t, equalDate)
}

func Test_Ads_GetByFilter_WithWrongFilters(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qwe@mail.ru")
	assert.NoError(t, err)

	_, err = client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	wrongFilters := queryParam{
		"published": "true",
		"user_id":   "1",
		"title":     "hello",
	}
	responseList, err := client.listAdsFilters(wrongFilters)
	assert.NoError(t, err)
	assert.Equal(t, len(responseList.Data), 0)

}

func Test_Ads_Delete(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qwerty@mail.ru")
	assert.NoError(t, err)

	ads, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	param := queryParam{"user_id": strconv.Itoa(int(user.Data.ID))}
	deleteAds, err := client.deleteAd(param, ads.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, deleteAds.AdId, ads.Data.ID)
	assert.Equal(t, deleteAds.AuthorId, user.Data.ID)

	deleteAds, err = client.deleteAd(param, ads.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, deleteAds.AdId, ads.Data.ID)
	assert.Equal(t, deleteAds.AuthorId, user.Data.ID)

}

func Test_Ads_Delete_Forbidden(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qwerty@mail.ru")
	assert.NoError(t, err)

	ads, err := client.createAd(user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	param := queryParam{"user_id": strconv.Itoa(int(user.Data.ID + 1))}
	_, err = client.deleteAd(param, ads.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}
