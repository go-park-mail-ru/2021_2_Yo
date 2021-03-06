package utils

import (
	log "backend/pkg/logger"
	"crypto/sha256"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/viper"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
)

const logMessage = "config:"

var (
	ErrFileExt = errors.New("wrong file extension")
	ErrFileDec = errors.New("file decoding error")
)

func CreatePasswordHash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func InitPostgresDB() (*sqlx.DB, error) {
	message := logMessage + "InitPostgresDB:"

	user := viper.GetString("postgres_db.user")
	password := os.Getenv("POSTGRES_PASSWORD")
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
	return db, nil
}

func InitRedisDB() (*redis.Client, error) {
	addr := viper.GetString("redis_db.addr")
	dbId := viper.GetInt("redis_db.db_id")
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       dbId,
	})
	if client == nil {
		return nil, errors.New("redis db init failed")
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
	case ".jpg",".jpeg":
	case ".png":
	case ".ico",".woff",".swg",".webp",".webm",".gif":
	default:
		return "", ErrFileExt
	}
	imgPath := viper.GetString("img_path")
	switch fileExtension {
	case ".jpg",".jpeg":
		newFileName := fileNameParts[0]+".webp"
		output, err := os.Create(imgPath+"/"+newFileName)
		if err != nil {
			log.Error(err)
		}
		defer output.Close()

		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 40)
		if err != nil {
			log.Error(err)
		}

		img, err := jpeg.Decode(file)
		if err != nil {
			return "", ErrFileDec
		}

		if err := webp.Encode(output, img, options); err != nil {
			log.Error(err)
		}

		return "https://bmstusa.ru/images/" + newFileName, nil

	case ".png":
		newFileName := fileNameParts[0]+".webp"
		output, err := os.Create(imgPath+"/"+newFileName)
		if err != nil {
			log.Error(err)
		}
		defer output.Close()

		options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 40)
		if err != nil {
			log.Error(err)
		}

		img, err := png.Decode(file)
		if err != nil {
			return "", ErrFileDec
		}

		if err := webp.Encode(output, img, options); err != nil {
			log.Error(err)
		}

		return "https://bmstusa.ru/images/" + newFileName, nil
		
	case ".ico",".woff",".swg",".webp",".webm",".gif":
	default:
		return "", ErrFileDec
	}

	output, err := os.Create(filepath.Join(imgPath, filepath.Base(fileName)))
	if err != nil {
		return "", err
	}
	defer output.Close()

	_, err = io.Copy(output, file)
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
