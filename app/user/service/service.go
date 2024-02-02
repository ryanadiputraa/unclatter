package service

import (
	"context"

	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/pkg/logger"
)

type service struct {
	log        logger.Logger
	repository user.UserRepository
}

func NewService(log logger.Logger, repository user.UserRepository) user.UserService {
	return &service{
		log:        log,
		repository: repository,
	}
}

func (s *service) CreateUser(ctx context.Context, arg user.NewUserArg) (*user.User, error) {
	user, err := user.NewUser(arg)
	if err != nil {
		s.log.Warn("user service: fail to create new user", err.Error())
		return nil, err
	}

	if err := s.repository.SaveOrUpdate(ctx, *user); err != nil {
		s.log.Error("user service: fail to save user", err.Error())
		return nil, err
	}

	return user, nil
}
