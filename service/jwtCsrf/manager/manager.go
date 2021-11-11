package middleware

import (
	"time"
	"github.com/dgrijalva/jwt-go/v4"
)

type jwtCsrfManager struct {
	secretWord []byte
}

func newjwtCsrfManager () *jwtCsrfManager {
	return &jwtCsrfManager{
		secretWord: []byte("secret1234"),
	}
}

func parseToken(susToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(susToken, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrHashUnavailable
		}
		return signingKey, nil
	})

	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		if (claims.ExpiresAt.Before(time.Now())) {
			return "", jwt.ErrHashUnavailable
		}
		return claims.ID, nil
	}
	return "", jwt.ErrHashUnavailable
}

func (m *jwtCsrfManager) Create(userId string) (string, error)  {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ID: userId,
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 7 * 24)), //Week  P.S. Maybe Frontend should ask us
	})

	csrfToken, err := jwtToken.SignedString(m.secretWord)

	return csrfToken, err
}

func (m *jwtCsrfManager) Check(susToken string) (string, error) {
	userId, err := parseToken(susToken, m.secretWord)
	if err != nil {
		return "", err
	}

	return userId, err
}