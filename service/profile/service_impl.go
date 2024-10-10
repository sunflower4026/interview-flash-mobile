package profile

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/sunflower4026/interview-flash-mobile/common/constants"
	userRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/user"
	profileWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/profile"
	userWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/user"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
	"gorm.io/gorm"
)

type ProfileServiceImpl struct {
	UserRepository userRepository.UserRepository
	DB             *gorm.DB
}

func NewProfileService(userRepository userRepository.UserRepository, DB *gorm.DB) ProfileService {
	return &ProfileServiceImpl{
		UserRepository: userRepository,
		DB:             DB,
	}
}

func (s *ProfileServiceImpl) Update(ctx context.Context, user profileWeb.UpdateProfileRequest) (userWeb.UserResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		log.FromCtx(ctx).Error(errors.New("user id not found in context"), "user id not found in context")
		return userWeb.UserResponse{}, errors.New("user id not found in context")
	}

	userToUpdate, err := s.UserRepository.Find(ctx, tx, fmt.Sprintf("id = '%s'", userID))
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find user")
		return userWeb.UserResponse{}, err
	}

	if user.FirstName != nil {
		userToUpdate.FirstName = *user.FirstName
	}
	if user.LastName != nil {
		userToUpdate.LastName = *user.LastName
	}
	if user.Address != nil {
		userToUpdate.Address = *user.Address
	}

	err = s.UserRepository.Update(ctx, tx, &userToUpdate)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to update user")
		return userWeb.UserResponse{}, err
	}

	tx.Commit()

	return userWeb.ToUserResponse(userToUpdate), nil
}
