package usecase

import (
	"backend/auth/repository/localstorage"
	"backend/models"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSignUp(t *testing.T) {
	secretWord := []byte("secret")
	repositoryMock := new(localstorage.RepositoryUserMock)
	useCaseTest := NewUseCaseAuth(repositoryMock, secretWord)

	userTest := &models.User{
		Name:     "nameTest",
		Surname:  "surnameTest",
		Mail:     "mailTest",
		Password: "passwordTest",
	}

	repositoryMock.On("CreateUser", userTest).Return(nil)
	err := useCaseTest.SignUp(userTest.Name, userTest.Surname, userTest.Mail, userTest.Password)
	require.NoError(t, err, err)
}

func TestSignIn(t *testing.T) {
	secretWord := []byte("secret")
	repositoryMock := new(localstorage.RepositoryUserMock)
	useCaseTest := NewUseCaseAuth(repositoryMock, secretWord)

	testId := "0"

	userTest := &models.User{
		ID: testId,
	}

	mail := "mailTest"
	password := "passwordTest"

	repositoryMock.On("GetUser", mail, password).Return(userTest, nil)
	signedStringTest, err := useCaseTest.SignIn(mail, password)
	require.NoError(t, err, "TestSignIn : useCase.SignIn err = ", err)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{ID: testId})
	signedString, err := token.SignedString(secretWord)
	require.NoError(t, err, "TestSignIn : token.SignedString err = ", err)

	require.Equal(t, signedString, signedStringTest, "TestSignIn : signedStrings are not equal")
}

func TestParseToken(t *testing.T) {
	secretWord := []byte("secret")
	repositoryMock := new(localstorage.RepositoryUserMock)
	useCaseTest := NewUseCaseAuth(repositoryMock, secretWord)

	userIdTest := "0"
	repositoryMock.On("GetUserById", userIdTest).Return(&models.User{
		ID: "0",
	}, nil)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims{ID: userIdTest})
	signedString, err := token.SignedString(useCaseTest.secretWord)
	require.NoError(t, err, "TestSignIn : token.SignedString err = ", err)
	result, err := useCaseTest.ParseToken(signedString)
	require.NoError(t, err, "TestSignIn : useCaseTest.ParseToken err = ", err)

	require.Equal(t, userIdTest, result.ID, "TestParseToken : expected and got IDs are not equal")
}
