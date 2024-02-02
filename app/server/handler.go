package server

import (
	authHandler "github.com/ryanadiputraa/unclatter/app/auth/handler"
	_userRepository "github.com/ryanadiputraa/unclatter/app/user/repository"
	_userService "github.com/ryanadiputraa/unclatter/app/user/service"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
)

func (s *Server) setupHandlers() {
	auth := s.web.Group("/auth")

	userRepository := _userRepository.NewRepository(s.db)
	userService := _userService.NewService(s.log, userRepository)

	googleOauth := oauth.NewGoogleOauth(s.config.GoogleOauth)

	authHandler.NewHandler(auth, s.config, s.log, userService, googleOauth)
}
