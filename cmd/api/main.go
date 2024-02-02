package main

import (
	"github.com/ryanadiputraa/unclatter/app/server"
	"github.com/ryanadiputraa/unclatter/config"
	"github.com/ryanadiputraa/unclatter/pkg/db/postgres"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
)

func main() {
	log := logger.NewLogger()

	config, err := config.LoadConfig("yml", "config/config.yml")
	if err != nil {
		log.Fatal("load config:", err)
	}

	db, err := postgres.NewDB(config)
	if err != nil {
		log.Fatal("db connection:", err)
	}

	server := server.NewHTTPServer(config, log, db)
	if err := server.ServeHTTP(); err != nil {
		log.Fatal("start server:", err)
	}
}
