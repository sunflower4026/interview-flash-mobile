package transaction

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	transactionWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/transaction"
	transactionService "gitlab.com/sunflower4026/interview-flash-mobile/service/transaction"
)

type TransactionControllerImpl struct {
	TransactionService transactionService.TransactionService
}

func NewTransactionController(transactionService transactionService.TransactionService) TransactionController {
	return &TransactionControllerImpl{
		TransactionService: transactionService,
	}
}

func (controller *TransactionControllerImpl) Topup(ctx *gin.Context) {
	var topupRequest transactionWeb.TopupRequest
	if error := ctx.ShouldBindJSON(&topupRequest); error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	transactionResponse, error := controller.TransactionService.Topup(ctx.Request.Context(), topupRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, transactionResponse, nil)
}

func (controller *TransactionControllerImpl) Payment(ctx *gin.Context) {
	var paymentRequest transactionWeb.PaymentRequest
	if error := ctx.ShouldBindJSON(&paymentRequest); error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	transactionResponse, error := controller.TransactionService.Payment(ctx.Request.Context(), paymentRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, transactionResponse, nil)
}

func (controller *TransactionControllerImpl) Transfer(ctx *gin.Context) {
	var transferRequest transactionWeb.TransferRequest
	if error := ctx.ShouldBindJSON(&transferRequest); error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	transactionResponse, error := controller.TransactionService.Transfer(ctx.Request.Context(), transferRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponseData(ctx, transactionResponse, nil)
}

func (controller *TransactionControllerImpl) FindAll(ctx *gin.Context) {
	filterRequest, err := httpservice.GetFilterRequest(ctx)
	if err != nil {
		httpservice.ResponsePagination(ctx, nil, httpservice.PaginationRequest{}, httpservice.PaginationResponse{}, err)
		return
	}

	paginationRequest, err := httpservice.GetPaginationRequest(ctx)
	if err != nil {
		httpservice.ResponsePagination(ctx, nil, httpservice.PaginationRequest{}, httpservice.PaginationResponse{}, err)
		return
	}

	orderRequest, err := httpservice.GetOrderRequest(ctx)
	if err != nil {
		httpservice.ResponsePagination(ctx, nil, paginationRequest, httpservice.PaginationResponse{}, err)
		return
	}

	transactionResponses, paginationResponse, error := controller.TransactionService.FindAll(ctx.Request.Context(), filterRequest, paginationRequest, orderRequest)
	if error != nil {
		httpservice.ResponseData(ctx, nil, error)
		return
	}

	httpservice.ResponsePagination(ctx, transactionResponses, paginationRequest, paginationResponse, nil)
}
