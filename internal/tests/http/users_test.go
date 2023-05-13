package http

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_User_Create(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Nickname, "qwertys")
	assert.Equal(t, response.Data.Email, "qw@mail.ru")

}

func Test_User_GetByID(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	userByID, err := client.getUserByID(user.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, userByID.Data.Nickname, "qwertys")
	assert.Equal(t, userByID.Data.Email, "qw@mail.ru")
}

func Test_User_GetByID_WrongID(t *testing.T) {
	client := getTestClient()

	_, err := client.getUserByID(1)
	assert.Error(t, err)
}

func Test_User_Update(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	userUpdate, err := client.updateUser(user.Data.ID, "qwertys1", "qw@mail.ru")
	assert.NoError(t, err)

	assert.Equal(t, userUpdate.Data.Nickname, "qwertys1")
	assert.Equal(t, userUpdate.Data.Email, "qw@mail.ru")
}

func Test_User_Delete(t *testing.T) {
	client := getTestClient()

	user, err := client.createUser("qwertys", "qw@mail.ru")
	assert.NoError(t, err)

	deleteUser, err := client.deleteUser(user.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, deleteUser.UserId, user.Data.ID)

	deleteUser, err = client.deleteUser(user.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, deleteUser.UserId, user.Data.ID)
}
