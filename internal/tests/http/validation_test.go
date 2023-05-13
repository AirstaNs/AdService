package http

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateAd_EmptyTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	_, err = client.createAd(userID, "", "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongTitle(t *testing.T) {
	client := getTestClient()

	title := strings.Repeat("a", 101)

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	_, err = client.createAd(userID, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_EmptyText(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	_, err = client.createAd(userID, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestCreateAd_TooLongText(t *testing.T) {
	client := getTestClient()

	text := strings.Repeat("a", 501)

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	_, err = client.createAd(userID, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(userID, resp.Data.ID, "", "new_world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongTitle(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	title := strings.Repeat("a", 101)

	_, err = client.updateAd(userID, resp.Data.ID, title, "world")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_EmptyText(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(userID, resp.Data.ID, "title", "")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestUpdateAd_TooLongText(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "@mail.ru")
	assert.NoError(t, err)
	userID := user.Data.ID

	text := strings.Repeat("a", 501)

	resp, err := client.createAd(userID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(userID, resp.Data.ID, "title", text)
	assert.ErrorIs(t, err, ErrBadRequest)
}
