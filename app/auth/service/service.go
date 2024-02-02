package service

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/auth"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
)

type service struct {
	log        logger.Logger
	repository auth.AuthProviderRepository
}

func NewService(log logger.Logger, repository auth.AuthProviderRepository) auth.AuthService {
	return &service{
		log:        log,
		repository: repository,
	}
}

func (s *service) AddUserAuthProvider(ctx context.Context, arg auth.NewAuthProviderArg) (*auth.AuthProvider, error) {
	provider := auth.NewAuthProvider(arg)
	if err := s.repository.Save(ctx, *provider); err != nil {
		s.log.Error("auth service: fail to save auth provider ", err.Error())
		return nil, err
	}

	return provider, nil
}
