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

func (s *service) CreateUser(ctx context.Context, arg user.NewUserArg) (created *user.User, err error) {
	created, err = user.NewUser(arg)
	if err != nil {
		s.log.Warn("user service: fail to create new user", err.Error())
		return
	}

	if err = s.repository.SaveOrUpdate(ctx, *created); err != nil {
		s.log.Error("user service: fail to save user", err.Error())
		return
	}
	return
}

func (s *service) GetUserInfo(ctx context.Context, userID string) (user *user.User, err error) {
	user, err = s.repository.FindByID(ctx, userID)
	if err != nil {
		s.log.Warn("user service: fail to get user by id", err.Error())
		return
	}
	return
}
