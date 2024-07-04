package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ryanadiputraa/unclatter/app/middleware"
	"github.com/ryanadiputraa/unclatter/config"
	_http "github.com/ryanadiputraa/unclatter/pkg/http"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
	"gorm.io/gorm"
)

type Server struct {
	config *config.Config
	log    logger.Logger
	web    *http.ServeMux
	db     *gorm.DB
	rw     _http.ResponseWriter
}

func NewHTTPServer(config *config.Config, log logger.Logger, db *gorm.DB) *Server {
	return &Server{
		config: config,
		log:    log,
		web:    http.NewServeMux(),
		db:     db,
		rw:     _http.NewResponseWriter(),
	}
}

func (s *Server) ServeHTTP() error {
	s.setupHandlers()
	handler := middleware.CORSMiddleware(s.web)
	handler = middleware.ThrottleMiddleware(handler, s.rw)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%v", s.config.Server.Port),
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	s.log.Info("received terminate, graceful shutdown", sig)

	tc, shutdown := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdown()

	return server.Shutdown(tc)
}
