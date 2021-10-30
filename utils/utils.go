package utils

import (
	log "backend/logger"
	"errors"
	"fmt"
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

func InitPostgresDB() (*sql.DB, error) {
	message := logMessage + "initDB:"
	log.Debug(message + "started")

	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetString("db.port")
	dbname := viper.GetString("db.dbname")
	sslmode := viper.GetString("db.sslmode")
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, user, dbname, password, sslmode)
	log.Debug(message+"connStr =", connStr)

	db, err := sql.Connect("postgres", connStr)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	log.Info("db status =", db.Stats())
	return db, nil
}
