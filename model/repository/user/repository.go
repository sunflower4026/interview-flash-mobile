package user

import (
	"context"

	userDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/user"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, tx *gorm.DB, user *userDomain.User) error
	Update(ctx context.Context, tx *gorm.DB, user *userDomain.User) error
	Find(ctx context.Context, tx *gorm.DB, whr string) (userDomain.User, error)
	FindAll(ctx context.Context, tx *gorm.DB) ([]userDomain.User, error)
}
