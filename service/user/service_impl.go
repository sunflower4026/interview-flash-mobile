package user

import (
	"context"

	userRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/user"
	userWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/user"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	UserRepository userRepository.UserRepository
	DB             *gorm.DB
}

func NewUserService(userRepository userRepository.UserRepository, DB *gorm.DB) UserService {
	return &UserServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
	}
}

func (s *UserServiceImpl) GetAll(ctx context.Context) ([]userWeb.UserResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	users, err := s.UserRepository.FindAll(ctx, tx)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find all users")
		return nil, err
	}

	userResponses := userWeb.ToUserResponses(users)

	return userResponses, nil
}
