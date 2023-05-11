package bootstrap

import (
	"fmt"
	"os"

	"github.com/iconmobile-dev/go-interview/config"
	"go.uber.org/zap"
)

// LoggerAndConfig is returning logger & config which are
// required for bootstrapping a service server
// is using CONFIG_FILE env var and if not set uses
// cfgFilePath. If cfgFilePath is not set then it tries to find the config
func LoggerAndConfig(serverName string, test bool) (*zap.SugaredLogger, config.Config) {
	var err error
	// init logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to init logger: %s", err.Error())
		os.Exit(1)
	}
	defer logger.Sync() // flushes buffer, if any
	log := logger.Sugar()

	var cfg *config.Config

	// load config
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		cfg, err = config.LoadDefaultConfig()
	} else {
		cfg, err = config.Load(configFile)
	}
	if err != nil {
		log.Errorw("config load error", "error", err)
		os.Exit(1)
	}

	// set service name
	cfg.Server.Name = serverName

	return log, *cfg
}
