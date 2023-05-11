package storage

import (
	"os"
	"testing"
)

var db *DB
var cache *Cache

func TestMain(m *testing.M) {
	// setup before tests
	var err error

	// bootstrap logger and config
	SetupLoggerAndConfig("storage", true)

	// database
	db, err = NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Errorw("error initializing postgres database", "error", err)
		os.Exit(1)
	}
	log.Infow("connected to Postgres", "host", cfg.DB.Host)

	err = db.Reset()
	if err != nil {
		log.Errorw("error resetting database", "error", err)
		os.Exit(1)
	}

	// cache
	cache, err = NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Errorw("error initializing redis database", "error", err)
		os.Exit(1)
	}
	log.Infow("connected to redis", "host", cfg.DB.Host)

	err = cache.Reset()
	if err != nil {
		log.Errorw("error resetting redis cache", "error", err)
		os.Exit(1)
	}

	// run tests
	code := m.Run()

	// shutdown after tests
	db.Close()

	os.Exit(code)
}
