package httpservice

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/hashing"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
)

type Response struct {
	Data interface{} `json:"data"`
	PaginationRequest
	PaginationResponse
	Message string `json:"message"`
}

type PaginationResponse struct {
	TotalPage int   `json:"total_page,omitempty"`
	TotalData int64 `json:"total_data,omitempty"`
}

type responseHashed struct {
	Data string `json:"data"`
}

// CheckMessageConstant checks if the message is a known constant
func CheckMessageConstant(message string) (int, string) {
	if code, ok := SuccessResponses[message]; ok {
		return code, message
	}
	if code, ok := ErrorResponses[message]; ok {
		return code, message
	}
	return http.StatusInternalServerError, message
}

func ResponseData(ctx *gin.Context, data interface{}, err error) {
	message := SUCCESS_OK
	statusCode := SuccessResponses[SUCCESS_OK]
	if err != nil {
		statusCode, message = CheckMessageConstant(err.Error())
	}
	ctx.JSON(statusCode, Response{Data: data, Message: message})
}

func ResponseDataPlain(ctx *gin.Context, data interface{}, err error) {
	message := SUCCESS_OK
	statusCode := SuccessResponses[SUCCESS_OK]
	if err != nil {
		statusCode, message = CheckMessageConstant(err.Error())
	}
	ctx.JSON(statusCode, gin.H{"data": data, "message": message})
}

func ResponseDataHashed(ctx *gin.Context, data interface{}, keyAes, ivAes string) {
	jsonString, err := json.Marshal(data)
	if err != nil {
		log.FromCtx(ctx.Request.Context()).Error(err, "failed to marshal response")
		ctx.JSON(ErrorResponses[ERR_INTERNAL_SERVER_ERROR], gin.H{"error": ERR_INTERNAL_SERVER_ERROR})
		return
	}

	plaintext, err := hashing.Encrypt(jsonString, []byte(keyAes), []byte(ivAes))
	if err != nil {
		log.FromCtx(ctx.Request.Context()).Error(err, "failed to encrypt response")
		ctx.JSON(ErrorResponses[ERR_INTERNAL_SERVER_ERROR], gin.H{"error": ERR_INTERNAL_SERVER_ERROR})
		return
	}

	response := &responseHashed{
		Data: plaintext,
	}

	ctx.JSON(SuccessResponses[SUCCESS_OK], response)
}

func ResponsePagination(ctx *gin.Context, data interface{}, requestPagination PaginationRequest, responsePagination PaginationResponse, err error) {
	message := SUCCESS_OK
	statusCode := SuccessResponses[SUCCESS_OK]
	if err != nil {
		statusCode, message = CheckMessageConstant(err.Error())
	}

	ctx.JSON(statusCode, Response{
		Data:               data,
		PaginationRequest:  requestPagination,
		PaginationResponse: responsePagination,
		Message:            message,
	})
}

func ResponseRaw(ctx *gin.Context, data interface{}) {
	ctx.Set("response", data)
	ctx.JSON(SuccessResponses[SUCCESS_OK], data)
}

func ResponseOKWithBlank(ctx *gin.Context) {
	ctx.JSON(SuccessResponses[SUCCESS_OK], "")
}
