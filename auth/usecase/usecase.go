package usecase

import (
	"backend/auth"
	log "backend/logger"
	"backend/models"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
)

const logMessage = "auth:usecase:usecase:"

type UseCase struct {
	repository auth.Repository
	secretWord []byte
}

func NewUseCase(userRepo auth.Repository, secretWord []byte) *UseCase {
	return &UseCase{
		repository: userRepo,
		secretWord: secretWord,
	}
}

/*

type TokenMetaInfo struct {
	AccessToken string
	RefreshToken string
	AccTokenExpires int64
	RefTokenExpires int64
}

type MegaClaims struct {

}

func(a* UseCase) CreateToken(ID int) (*TokenMetaInfo, error){
	var err error

	TokenMeta := &TokenMetaInfo{}
	TokenMeta.AccTokenExpires = time.Now().Add(time.Minute * 15).Unix()
	TokenMeta.RefTokenExpires = time.Now().Add(time.Hour * 7).Unix()

	//Делаю access token
	AccTokenClaims := jwt.MapClaims{}
	AccTokenClaims["authorized"] = true
	AccTokenClaims["ID"] = ID
	AccTokenClaims["exp"] = TokenMeta.RefTokenExpires

	AccToken := jwt.NewWithClaims(jwt.SigningMethodHS256, AccTokenClaims)
	TokenMeta.AccessToken, err = AccToken.SignedString(a.secretWord)

	if err != nil {
		return nil, err
	}
	
	//Делаю refresh token
	RefTokenClaims := jwt.MapClaims{}
	RefTokenClaims["ID"] = ID
	RefTokenClaims["exp"] = TokenMeta.RefTokenExpires

	RefToken := jwt.NewWithClaims(jwt.SigningMethodHS256, RefTokenClaims)
	TokenMeta.RefreshToken, err = RefToken.SignedString(a.secretWord)

	if err != nil {
		return nil, err
	}

	return TokenMeta, nil
}
*/

type claims struct {
	jwt.StandardClaims
	ID string `json:"user_id"`
	ExpiresAt int64 `json:"exp"`
}

func parseToken(accessToken string, signingKey []byte) (string, error) {
	token, err := jwt.ParseWithClaims(accessToken, &claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*claims); ok && token.Valid {
		return claims.ID, nil
	}
	return "", auth.ErrInvalidAccessToken
}

func (a *UseCase) SignUp(name, surname, mail, password string) error {
	password_hash := a.CreatePasswordHash(password)
	user := &models.User{
		Name:     name,
		Surname:  surname,
		Mail:     mail,
		Password: password_hash,
	}
	return a.repository.CreateUser(user)
}

func (a *UseCase) SignIn(mail, password string) (string, error) {
	message := logMessage + "SignIn:"
	hashedPassword := a.CreatePasswordHash(password)
	log.Debug(message+"hashedPassword =", hashedPassword)
	user, err := a.repository.GetUser(mail, hashedPassword)
	if err == auth.ErrUserNotFound {
		log.Error(message+"err =", err)
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims{
		ID: user.ID,
		ExpiresAt: time.Now().Add(time.Minute * 30).Unix(),
	})
	return token.SignedString(a.secretWord)
}

func (a *UseCase) ParseToken(accessToken string) (*models.User, error) {
	userID, err := parseToken(accessToken, a.secretWord)
	if err != nil {
		return nil, err
	}
	return a.repository.GetUserById(userID)
}

func (a *UseCase) CreatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
