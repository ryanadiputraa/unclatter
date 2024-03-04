package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ryanadiputraa/unclatter/app/middleware"
	"github.com/ryanadiputraa/unclatter/config"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	log    logger.Logger
	web    *http.ServeMux
	db     *gorm.DB
}

func NewHTTPServer(config *config.Config, log logger.Logger, db *gorm.DB) *Server {
	return &Server{
		config: config,
		log:    log,
		web:    http.NewServeMux(),
		db:     db,
	}
}

func (s *Server) ServeHTTP() error {
	s.setupHandlers()
	handler := middleware.CORSMiddleware(s.web)

	server := &http.Server{
		Addr:         s.config.Server.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 30,
		WriteTimeout: time.Second * 30,
	}

	go func() {
		s.log.Info("starting server on port", s.config.Server.Port)
		if err := server.ListenAndServe(); err != nil {
			s.log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, os.Kill)

	sig := <-quit
	s.log.Info("received terminate, graceful shutdown", sig)

	tc, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdown()

	return server.Shutdown(tc)
}
