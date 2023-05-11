// Package userlib - Authentication logic
package userlib

import (
	"os"
	"testing"

	"github.com/iconmobile-dev/go-interview/lib/storage"
	"github.com/jmoiron/sqlx"
)

var db *storage.DB
var failingDB *storage.DB
var cache *storage.Cache

func TestMain(m *testing.M) {
	// setup before tests
	var err error

	// bootstrap logger and config
	SetupLoggerAndConfig("userlib", true)

	// database
	db, err = storage.NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
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
	cache, err = storage.NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Errorw("error initializing redis cache database", "error", err)
		os.Exit(1)
	}
	log.Infow("connected to redis", "host", cfg.DB.Host)

	err = cache.Reset()
	if err != nil {
		log.Errorw("error resetting cache", "error", err)
		os.Exit(1)
	}

	failingDB = &storage.DB{}
	{
		db, err := sqlx.Open("postgres", "")
		if err != nil {
			log.Errorw("error resetting database", "error", err)
			os.Exit(1)
		}
		failingDB.DB = db
	}

	// run tests
	code := m.Run()

	// shutdown after tests
	db.Close()
	cache.Close()

	os.Exit(code)
}
