package localstorage

import (
	"backend/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateUser(t *testing.T) {
	repositotyTest := NewRepository()
	err := repositotyTest.CreateUser(&models.User{})
	require.NoError(t, err, "TestCreateUser : err = ", err)
}

func TestGetUser(t *testing.T) {
	repositotyTest := NewRepository()
	idInt := 0
	mail := "mailTest"
	password := "passwordTest"
	userToAppend := &User{
		ID:       idInt,
		Mail:     mail,
		Password: password,
	}
	userExpected := toModelUser(userToAppend)
	repositotyTest.users = append(repositotyTest.users, userToAppend)
	userTest, err := repositotyTest.GetUser(mail, password)
	require.NoError(t, err, "TestGetUser : repository.GetUser err = ", err)
	require.Equal(t, userExpected, userTest, "TestGetUser : expected and got users are not equal")
}

func TestGetUserById(t *testing.T) {
	repositotyTest := NewRepository()
	idInt := 0
	idString := "0"
	userToAppend := &User{
		ID: idInt,
	}
	repositotyTest.users = append(repositotyTest.users, userToAppend)
	userTest, err := repositotyTest.GetUserById(idString)
	require.NoError(t, err, "TestGetUserById : repository.GetUserById err = ", err)
	userExpected := toModelUser(userToAppend)
	require.Equal(t, userExpected, userTest, "TestGetUserById : expected and got users are not equal")
}
