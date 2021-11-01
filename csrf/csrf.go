package csrf

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var Token *HashToken

type HashToken struct {
	CSRFSecretWord []byte
}

func NewHMACHashToken(secret string) (*HashToken, error) {
	return &HashToken{CSRFSecretWord: []byte(secret)}, nil
}

func (token *HashToken) Create(cookie string, expirationTime int64) (string, error) {
	hash := hmac.New(sha256.New, token.CSRFSecretWord)
	data := fmt.Sprintf("%s:%d", cookie, expirationTime)
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", err
	}
	CSRFToken := hex.EncodeToString(hash.Sum(nil)) + ":" + strconv.FormatInt(expirationTime, 10)
	return CSRFToken, nil
}

func (token *HashToken) Check(cookie string, CSRFToken string) (bool, error) {
	tokenData := strings.Split(CSRFToken, ":")
	if len(tokenData) != 2 {
		return false, ErrBadToken
	}

	tokenExpiration, err := strconv.ParseInt(tokenData[1], 10, 64)
	if err != nil {
		return false, ErrTokenNoExp
	}

	if tokenExpiration < time.Now().Unix() {
		return false, ErrTokenExp
	}

	hash := hmac.New(sha256.New, token.CSRFSecretWord)
	data := fmt.Sprintf("%s:%d", cookie, tokenExpiration)
	_, err = hash.Write([]byte(data))
	if err != nil {
		return false, err
	}
	expectedCSRFToken := hash.Sum(nil)
	gottenCSRFToken, err := hex.DecodeString(tokenData[0])
	if err != nil {
		return false, err
	}
	return hmac.Equal(gottenCSRFToken, expectedCSRFToken), nil
}
