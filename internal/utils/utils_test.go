package utils

import (
	"testing"
)

func TestCreatePasswordHash(t *testing.T) {
	_ = CreatePasswordHash("test")
}

func TestInitPostgresDB(t *testing.T) {
	_, _ = InitPostgresDB()
}

func TestInitRedisDB(t *testing.T) {
	_, _ = InitRedisDB()
}
