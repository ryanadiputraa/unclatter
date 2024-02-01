package service

import (
	"fmt"

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

func (s *service) CreateUser(arg user.CreateUserArg) (*user.User, error) {
	user, err := user.NewUser(arg)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Save(*user); err != nil {
		s.log.Info(fmt.Sprintf("user service: fail to save user \"%v\"", err.Error()))
		return nil, err
	}

	s.log.Info(fmt.Sprintf("user service: register new user \"%v\"", user.ID))
	return user, nil
}
