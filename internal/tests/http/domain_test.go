package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(100, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(100, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(userID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(userID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}

func TestUpdateAd_changedUpdateTime(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)
	time.Sleep(time.Second)
	_, err = client.updateAd(userID, resp.Data.ID, "title", "text")
	assert.NoError(t, err)

	resp, err = client.getAdByID(resp.Data.ID)
	assert.NoError(t, err)
	assert.NotEqual(t, resp.Data.CreateDate, resp.Data.UpdateDate)
}
