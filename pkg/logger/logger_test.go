package logger

import (
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLogger(t *testing.T) {
	Init(log.DebugLevel)
	Debug("test")
	Info("test")
	Error("test")
}
