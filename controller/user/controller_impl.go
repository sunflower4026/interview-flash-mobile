package user

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	userService "gitlab.com/sunflower4026/interview-flash-mobile/service/user"
)

type UserControllerImpl struct {
	UserService userService.UserService
}

func NewUserController(userService userService.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) FindAll(ctx *gin.Context) {
	userResponses, error := controller.UserService.GetAll(ctx.Request.Context())
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, userResponses, nil)
}
