package profile

import (
	"github.com/gin-gonic/gin"
)

type ProfileController interface {
	Update(ctx *gin.Context)
}
