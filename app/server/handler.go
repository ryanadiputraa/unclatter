package server

import (
	authHandler "github.com/ryanadiputraa/unclatter/app/auth/handler"
	_authRepository "github.com/ryanadiputraa/unclatter/app/auth/repository"
	_authService "github.com/ryanadiputraa/unclatter/app/auth/service"
	"github.com/ryanadiputraa/unclatter/app/middleware"
	userHandler "github.com/ryanadiputraa/unclatter/app/user/handler"
	_userRepository "github.com/ryanadiputraa/unclatter/app/user/repository"
	_userService "github.com/ryanadiputraa/unclatter/app/user/service"
	_http "github.com/ryanadiputraa/unclatter/pkg/http"
	"github.com/ryanadiputraa/unclatter/pkg/jwt"
	"github.com/ryanadiputraa/unclatter/pkg/oauth"
)

func (s *Server) setupHandlers() {
	// e := echo.New()
	// s.web.Handle("/", middleware.CORSMiddleware())

	// auth := s.web.Group("/auth")
	// article := s.web.Group("/api/articles")

	rw := _http.NewResponseWriter()
	// validator := validator.NewValidator()
	googleOauth := oauth.NewGoogleOauth(s.config.GoogleOauth)
	jwtTokens := jwt.NewJWTTokens(s.config.JWT)
	// scrapper := scrapper.NewScrapper()
	// sanitizer := sanitizer.NewSanitizer()

	authMiddleware := middleware.NewAuthMiddleware(s.log, s.config.JWT, rw, jwtTokens)

	userRepository := _userRepository.NewRepository(s.db)
	userService := _userService.NewService(s.log, userRepository)
	userHandler.NewUserHandler(s.web, rw, userService, *authMiddleware)

	authRepository := _authRepository.NewRepository(s.db)
	authService := _authService.NewService(s.log, authRepository)
	authHandler.NewHandler(s.web, s.config, s.log, authService, userService, googleOauth, jwtTokens)

	// articleRepository := _articleRepository.NewRepository(s.db)
	// articleService := _articleService.NewService(s.log, scrapper, sanitizer, articleRepository)
	// articleHandler.NewHandler(article, articleService, *authMiddleware, validator)
}
