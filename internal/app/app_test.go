package app

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApp(t *testing.T) {
	options := &Options{
		LogLevel: log.DebugLevel,
		Testing:  true,
	}
	app, err := NewApp(options)
	require.NoError(t, err)
	err = app.Run()
	require.Error(t, err)
}
