package transaction

import (
	"github.com/gin-gonic/gin"
)

type TransactionController interface {
	Topup(ctx *gin.Context)
	Payment(ctx *gin.Context)
	Transfer(ctx *gin.Context)

	FindAll(ctx *gin.Context)
}
