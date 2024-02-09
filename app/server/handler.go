package server

import (
	articleHandler "github.com/ryanadiputraa/unclatter/app/article/handler"
	_articleService "github.com/ryanadiputraa/unclatter/app/article/service"
	authHandler "github.com/ryanadiputraa/unclatter/app/auth/handler"
	_authRepository "github.com/ryanadiputraa/unclatter/app/auth/repository"
	_authService "github.com/ryanadiputraa/unclatter/app/auth/service"
	"github.com/ryanadiputraa/unclatter/app/middleware"
	userHandler "github.com/ryanadiputraa/unclatter/app/user/handler"
	_userRepository "github.com/ryanadiputraa/unclatter/app/user/repository"
	_userService "github.com/ryanadiputraa/unclatter/app/user/service"
	"github.com/ryanadiputraa/unclatter/pkg/jwt"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
)

func (s *Server) setupHandlers() {
	auth := s.web.Group("/auth")
	user := s.web.Group("/api/users")
	article := s.web.Group("/api/articles")

	googleOauth := oauth.NewGoogleOauth(s.config.GoogleOauth)
	jwtTokens := jwt.NewJWTTokens(s.config.JWT)
	scrapperLib := scrapper.NewScrapper()

	authMiddleware := middleware.NewAuthMiddleware(s.log, s.config.JWT, jwtTokens)

	userRepository := _userRepository.NewRepository(s.db)
	userService := _userService.NewService(s.log, userRepository)
	userHandler.NewUserHandler(user, userService, *authMiddleware)

	authRepository := _authRepository.NewRepository(s.db)
	authService := _authService.NewService(s.log, authRepository)
	authHandler.NewHandler(auth, s.config, s.log, authService, userService, googleOauth, jwtTokens)

	articleService := _articleService.NewService(s.log, scrapperLib)
	articleHandler.NewHandler(article, articleService, *authMiddleware)
}
