package usecase

import (
	protoAuth "backend/microservice/auth/proto"
	"backend/pkg/utils"
	"context"
	"github.com/dgrijalva/jwt-go/v4"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
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
		if claims.ExpiresAt.Before(time.Now()) {
			return "", jwt.ErrHashUnavailable
		}
		return claims.ID, nil
	}
	return "", jwt.ErrHashUnavailable
}

func (s *authService) CreateToken(ctx context.Context, protoUserId *protoAuth.UserId) (*protoAuth.CSRFToken, error) {
	message := logMessage + "CreateToken:"
	log.Debug(message + "started")
	csrfToken, err := utils.GenerateCsrfToken(protoUserId.ID)
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
		ID: userId,
	}
	return response, err
}
