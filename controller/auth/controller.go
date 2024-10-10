package auth

import (
	"github.com/gin-gonic/gin"
)

type AuthController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}
