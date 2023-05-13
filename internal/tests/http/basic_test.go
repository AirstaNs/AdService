package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Добавлено создание пользователя. (ID 123 нет в базе) Тк включена проверка на существующего пользователя. Чтобы удовлетворить условию:
"в существующие методы необходимо добавить проверку, что пользователь с переданным ID существует,
а именно нужно получать AuthorID из места хранения"
*/

func TestCreateAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	response, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, userID)
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	response, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(userID, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(userID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(userID, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	response, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(userID, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	response, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(userID, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(userID, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}
