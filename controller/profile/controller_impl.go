package profile

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	profileWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/profile"
	profileService "gitlab.com/sunflower4026/interview-flash-mobile/service/profile"
)

type ProfileControllerImpl struct {
	ProfileService profileService.ProfileService
}

func NewProfileController(profileService profileService.ProfileService) ProfileController {
	return &ProfileControllerImpl{
		ProfileService: profileService,
	}
}

func (controller *ProfileControllerImpl) Update(ctx *gin.Context) {
	var updateProfileRequest profileWeb.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&updateProfileRequest); err != nil {
		httpservice.ResponseData(ctx, nil, err)
		return
	}

	userResponse, error := controller.ProfileService.Update(ctx.Request.Context(), updateProfileRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, userResponse, nil)
}
