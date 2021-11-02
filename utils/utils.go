package utils

import (
	log "backend/logger"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"github.com/gomodule/redigo/redis"
	sql "github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"os"
)

const logMessage = "config:"

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
