// Package main contains the binary
//
package main

import (
	"fmt"
	"os"

	"github.com/iconmobile-dev/go-interview/lib/bootstrap"
	"github.com/iconmobile-dev/go-interview/lib/storage"
	"github.com/iconmobile-dev/go-interview/services/user"
)

func main() {
	// bootstrap logger and config
	log, cfg := bootstrap.LoggerAndConfig("user", false)

	// open database
	db, err := storage.NewDB(cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	if err != nil {
		log.Errorw("error initializing database postgres", "error", err)
		os.Exit(1)
	}

	// open cache
	cache, err := storage.NewCache(cfg.Redis.Host, cfg.Redis.Port, cfg.Redis.Password)
	if err != nil {
		log.Errorw("error initializing cache database redis new", "error", err)
		os.Exit(1)
	}

	// init service
	s := user.New(db, cache)

	log.Infow("Starting", cfg.Server.Name, "on", cfg.Server.Env, "using port", cfg.Server.PortEngagement)

	bind := fmt.Sprintf(":%v", cfg.Server.PortEngagement)
	srv := bootstrap.Server(s, bind)

	err = srv.ListenAndServe()
	if err != nil {
		log.Errorw("error starting server", "error", err)
		os.Exit(3)
	}
}
