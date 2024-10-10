package user

import (
	"context"

	userWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/user"
)

type UserService interface {
	GetAll(ctx context.Context) ([]userWeb.UserResponse, error)
}
