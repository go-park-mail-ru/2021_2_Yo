package utils

import (
	log "backend/logger"
	"backend/response"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	sql "github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const logMessage = "config:"

var (
	ErrFileExt = errors.New("wrong file extension")
)

func GetSecret() (string, error) {
	message := logMessage + "getSecret:"
	log.Debug(message + "started")
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret1234"
		err := errors.New("Can't get secret from environment")
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

func InitPostgresDB() (*sql.DB, error) {
	message := logMessage + "InitPostgresDB:"
	log.Debug(message + "started")

	user := viper.GetString("postgres_db.user")
	password := viper.GetString("postgres_db.password")
	host := viper.GetString("postgres_db.host")
	port := viper.GetString("postgres_db.port")
	dbname := viper.GetString("postgres_db.dbname")
	sslmode := viper.GetString("postgres_db.sslmode")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	log.Debug(message+"connStr =", connStr)

	db, err := sql.Connect("postgres", connStr)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	log.Info(message+"db status =", db.Stats())
	return db, nil
}

func InitRedisDB(dbConfName string) (redis.Conn, error) {
	message := logMessage + "InitRedisDB:"
	log.Debug(message + "started")

	name := viper.GetString(dbConfName + ".name")
	value := viper.GetString(dbConfName + ".value")
	usage := viper.GetString(dbConfName + ".usage")
	log.Debug(message+"name,value,usage =", name, value, usage)

	redisAddr := flag.String(name, value, usage)
	return redis.DialURL(*redisAddr)
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
	fileExtension := filepath.Ext(fileName)
	switch fileExtension {
	case ".jpg":
	case ".jpeg":
	case ".png":
	case ".ico":
	case ".woff":
	case ".swg":
	case ".webp":
	case ".webm":
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
	log.Debug(message+"imgUrl =", "https://bmstusa.ru/images/"+fileName)
	return "https://bmstusa.ru/images/" + fileName, nil
}

func CheckIfNoError(w *http.ResponseWriter, err error, msg string, status response.HttpStatus) bool {
	if err != nil {
		log.Error(msg+"err =", err)
		response.SendResponse(*w, response.ErrorResponse(err.Error()))
		return false
	}
	return true
}
