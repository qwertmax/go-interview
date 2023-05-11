package storage

import (
	"github.com/iconmobile-dev/go-interview/config"
	"github.com/iconmobile-dev/go-interview/lib/bootstrap"
	"go.uber.org/zap"
)

var log *zap.SugaredLogger
var cfg config.Config

// SetupLoggerAndConfig sets the global logger and config dependency
// should be called during tests
func SetupLoggerAndConfig(serverName string, test bool) {
	log, cfg = bootstrap.LoggerAndConfig(serverName, test)
}

// initiates log and cfg with default values
func init() {
	SetupLoggerAndConfig("media", false)
}
