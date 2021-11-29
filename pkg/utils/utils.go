package utils

import (
	log "backend/pkg/logger"
	"backend/pkg/response"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const logMessage = "config:"

var (
	ErrFileExt = errors.New("wrong file extension")
)

func GetSecret() (string, error) {
	message := logMessage + "getSecret:"
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret1234"
		err := errors.New("can't get secret from environment")
		log.Error(message+"err =", err)
		return secret, nil
	}
	return secret, nil
}

func CreatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func InitPostgresDB() (*sqlx.DB, error) {
	message := logMessage + "InitPostgresDB:"

	user := viper.GetString("postgres_db.user")
	password := viper.GetString("postgres_db.password")
	host := viper.GetString("postgres_db.host")
	port := viper.GetString("postgres_db.port")
	dbname := viper.GetString("postgres_db.dbname")
	sslmode := viper.GetString("postgres_db.sslmode")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	log.Info(message+"db status =", db.Stats())
	return db, nil
}

func InitRedisDB(dbConfName string) (*redis.Client, error) {

	addr := viper.GetString(dbConfName + ".addr")
	dbId := viper.GetInt(dbConfName + ".db_id")

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       dbId,
	})
	if client == nil {
		return nil, errors.New("init redis db")
	}
	return client, nil
}

func SaveImageFromRequest(r *http.Request, key string) (string, error) {
	message := logMessage + "SaveImageFromRequest"
	_ = message
	file, handler, err := r.FormFile(key)
	if err != nil {
		return "", err
	}
	defer file.Close()
	imgUuid := uuid.NewV4()
	fileNameParts := strings.Split(handler.Filename, ".")
	fileNameParts[0] = imgUuid.String()
	fileName := fileNameParts[0] + "." + fileNameParts[1]
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	switch fileExtension {
	case ".jpg":
	case ".jpeg":
	case ".png":
	case ".ico":
	case ".woff":
	case ".swg":
	case ".webp":
	case ".webm":
	case ".gif":
	default:
		return "", ErrFileExt
	}
	dst, err := os.Create(filepath.Join("/home/ubuntu/go/2021_2_Yo/static/images", filepath.Base(fileName)))
	if err != nil {
		return "", err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}
	return "https://bmstusa.ru/images/" + fileName, nil
}

func GenerateCsrfToken(userId string) (string, error) {
	message := logMessage + "GenerateCsrfToken:"
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		ID:        userId,
		ExpiresAt: jwt.At(time.Now().Add(time.Hour * 7 * 24)), //Week  P.S. Maybe Frontend should ask us
	})
	secretWord := os.Getenv("CSRFSECRET")
	csrfToken, err := jwtToken.SignedString([]byte(secretWord))
	if err != nil {
		log.Error(message+"err = ", err)
		return "", err
	}
	return csrfToken, err
}

func CheckIfNoError(w *http.ResponseWriter, err error, msg string, status response.HttpStatus) bool {
	if err != nil {
		log.Error(msg+"err =", err)
		response.SendResponse(*w, response.ErrorResponse(err.Error()))
		return false
	}
	return true
}

type modifiedResponse struct {
	http.ResponseWriter
	StatusCode int
}

func NewModifiedResponse(w http.ResponseWriter) *modifiedResponse {
	return &modifiedResponse{ResponseWriter: w}
}

func (w *modifiedResponse) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
