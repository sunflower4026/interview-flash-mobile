package user

import (
	"context"

	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"

	userDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/user"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (r *UserRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, user *userDomain.User) error {
	if err := tx.Create(user).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create user")
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, tx *gorm.DB, user *userDomain.User) error {
	if err := tx.Save(user).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to update user")
		return err
	}

	return nil
}

func (r *UserRepositoryImpl) Find(ctx context.Context, tx *gorm.DB, where string) (userDomain.User, error) {
	var user userDomain.User
	if err := tx.Where(where).First(&user).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find user")
		return user, err
	}

	return user, nil
}

func (r *UserRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB) ([]userDomain.User, error) {
	var users []userDomain.User
	if err := tx.Find(&users).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find all users")
		return nil, err
	}

	return users, nil
}
