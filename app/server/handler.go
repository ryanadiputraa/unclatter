package server

import (
	authHandler "github.com/ryanadiputraa/unclatter/app/auth/handler"
	_authRepository "github.com/ryanadiputraa/unclatter/app/auth/repository"
	_authService "github.com/ryanadiputraa/unclatter/app/auth/service"
	userHandler "github.com/ryanadiputraa/unclatter/app/user/handler"
	_userRepository "github.com/ryanadiputraa/unclatter/app/user/repository"
	_userService "github.com/ryanadiputraa/unclatter/app/user/service"
	"github.com/ryanadiputraa/unclatter/pkg/jwt"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
)

func (s *Server) setupHandlers() {
	auth := s.web.Group("/auth")
	user := s.web.Group("/api/users")

	googleOauth := oauth.NewGoogleOauth(s.config.GoogleOauth)
	jwtTokens := jwt.NewJWTTokens(s.config.JWT)

	userRepository := _userRepository.NewRepository(s.db)
	userService := _userService.NewService(s.log, userRepository)
	userHandler.NewUserHandler(user, userService)

	authRepository := _authRepository.NewRepository(s.db)
	authService := _authService.NewService(s.log, authRepository)
	authHandler.NewHandler(auth, s.config, s.log, authService, userService, googleOauth, jwtTokens)
}
