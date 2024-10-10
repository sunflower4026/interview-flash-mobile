package auth

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	authWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/auth"
	authService "gitlab.com/sunflower4026/interview-flash-mobile/service/auth"
)

type AuthControllerImpl struct {
	AuthService authService.AuthService
}

func NewAuthController(authService authService.AuthService) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
	}
}

func (controller *AuthControllerImpl) Register(ctx *gin.Context) {
	var registerRequest authWeb.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		httpservice.ResponseData(ctx, nil, err)
		return
	}

	registerResponse, error := controller.AuthService.Register(ctx.Request.Context(), registerRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, registerResponse, nil)
}

func (controller *AuthControllerImpl) Login(ctx *gin.Context) {
	var loginRequest authWeb.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		httpservice.ResponseData(ctx, nil, err)
		return
	}

	loginResponse, error := controller.AuthService.Login(ctx.Request.Context(), loginRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, loginResponse, nil)
}
