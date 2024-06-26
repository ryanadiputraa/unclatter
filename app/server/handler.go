package server

import (
	"net/http"

	articleHandler "github.com/ryanadiputraa/unclatter/app/article/handler"
	_articleRepository "github.com/ryanadiputraa/unclatter/app/article/repository"
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
	"github.com/ryanadiputraa/unclatter/pkg/sanitizer"
	"github.com/ryanadiputraa/unclatter/pkg/scrapper"
	"github.com/ryanadiputraa/unclatter/pkg/validator"
)

func (s *Server) setupHandlers() {
	validator := validator.NewValidator()
	googleOauth := oauth.NewGoogleOauth(s.config.GoogleOauth)
	jwtTokens := jwt.NewJWTTokens(s.config.JWT)
	scrapper := scrapper.NewScrapper()
	sanitizer := sanitizer.NewSanitizer()

	authMiddleware := middleware.NewAuthMiddleware(s.log, s.config.JWT, s.rw, jwtTokens)

	userRepository := _userRepository.NewRepository(s.db)
	userService := _userService.NewService(s.log, userRepository)
	userHandler.NewUserHandler(s.web, s.rw, userService, *authMiddleware)

	authRepository := _authRepository.NewRepository(s.db)
	authService := _authService.NewService(s.log, authRepository)
	authHandler.NewHandler(s.web, s.config, s.log, authService, userService, googleOauth, jwtTokens)

	articleRepository := _articleRepository.NewRepository(s.db)
	articleService := _articleService.NewService(s.log, scrapper, sanitizer, articleRepository)
	articleHandler.NewHandler(s.web, s.rw, articleService, *authMiddleware, validator)

	s.web.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		s.rw.WriteResponseData(w, 200, "ok")
	})
}
