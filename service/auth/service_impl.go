package auth

import (
	"context"
	"fmt"

	"gitlab.com/sunflower4026/interview-flash-mobile/helper"
	userDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/user"
	userRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/user"
	authWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/auth"
	jwtService "gitlab.com/sunflower4026/interview-flash-mobile/service/jwt"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"

	"gorm.io/gorm"
)

type AuthServiceImpl struct {
	UserRepository userRepository.UserRepository
	JWTService     jwtService.JWTService
	DB             *gorm.DB
}

func NewAuthService(userRepository userRepository.UserRepository, JWTService jwtService.JWTService, DB *gorm.DB) AuthService {
	return &AuthServiceImpl{
		UserRepository: userRepository,
		JWTService:     JWTService,
		DB:             DB,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, user authWeb.RegisterRequest) (authWeb.RegisterResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	hashedPin, err := helper.GenerateHash(user.Pin)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to generate hash")
		return authWeb.RegisterResponse{}, err
	}

	userDomain := userDomain.User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		PhoneNumber: user.PhoneNumber,
		Pin:         hashedPin,
		Address:     user.Address,
	}
	if err := s.UserRepository.Create(ctx, tx, &userDomain); err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create user")
		return authWeb.RegisterResponse{}, err
	}

	tx.Commit()

	return authWeb.ToRegisterResponse(userDomain), nil
}

func (s *AuthServiceImpl) Login(ctx context.Context, userRequest authWeb.LoginRequest) (authWeb.LoginResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	user, err := s.UserRepository.Find(ctx, tx, fmt.Sprintf("phone_number = '%s'", userRequest.PhoneNumber))
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find user")
		return authWeb.LoginResponse{}, err
	}

	tx.Commit()

	if !helper.ValidatePassword(userRequest.Pin, user.Pin) {
		return authWeb.LoginResponse{}, fmt.Errorf("phone number and pin do not match")
	}

	accessToken, err := s.JWTService.GenerateAccessToken(user.ID.String())
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to generate access token")
		return authWeb.LoginResponse{}, err
	}

	refreshToken, err := s.JWTService.GenerateRefreshToken(user.ID.String())
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to generate refresh token")
		return authWeb.LoginResponse{}, err
	}

	return authWeb.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
