package auth

import (
	"context"

	authWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/auth"
)

type AuthService interface {
	Register(ctx context.Context, user authWeb.RegisterRequest) (authWeb.RegisterResponse, error)
	Login(ctx context.Context, user authWeb.LoginRequest) (authWeb.LoginResponse, error)
}
