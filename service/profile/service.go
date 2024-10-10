package profile

import (
	"context"

	profileWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/profile"
	userWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/user"
)

type ProfileService interface {
	Update(ctx context.Context, user profileWeb.UpdateProfileRequest) (userWeb.UserResponse, error)
}
