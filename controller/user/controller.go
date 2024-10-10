package user

import (
	"github.com/gin-gonic/gin"
)

type UserController interface {
	FindAll(ctx *gin.Context)
}
