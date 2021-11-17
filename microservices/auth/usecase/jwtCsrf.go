package usecase

import (
	"time"
	"github.com/dgrijalva/jwt-go/v4"
	"os"
	"context"
	protoAuth "backend/microservices/proto/auth"
	log "github.com/sirupsen/logrus"
)

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

func (s *authService) CreateToken(ctx context.Context, protoUserId *protoAuth.UserId) (*protoAuth.CSRFToken, error)  {
	message := logMessage + "CreateToken:"
	log.Debug(message + "started")
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ID: protoUserId.UserId,
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 7 * 24)), //Week  P.S. Maybe Frontend should ask us
	})
	secretWord := os.Getenv("CSRFSECRET")
	csrfToken, err := jwtToken.SignedString([]byte(secretWord))
	if err != nil {
		return &protoAuth.CSRFToken{}, err
	}
	response := &protoAuth.CSRFToken{
		CSRFToken: csrfToken,
	}
	return response, err
}

func (s *authService) CheckToken(ctx context.Context, protoToken *protoAuth.CSRFToken) (*protoAuth.UserId, error) {
	message := logMessage + "CheckToken:"
	log.Debug(message + "started")
	secretWord := os.Getenv("CSRFSECRET")
	susToken := protoToken.CSRFToken
	userId, err := parseToken(susToken, []byte(secretWord))
	if err != nil {
		return &protoAuth.UserId{}, err
	}
	response := &protoAuth.UserId{
		UserId: userId,
	}
	return response, err
}