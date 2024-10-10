package transaction

import (
	"context"

	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	transactionWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/transaction"
)

type TransactionService interface {
	Topup(ctx context.Context, topupRequest transactionWeb.TopupRequest) (transactionWeb.TransactionResponse, error)
	Payment(ctx context.Context, paymentRequest transactionWeb.PaymentRequest) (transactionWeb.TransactionResponse, error)
	Transfer(ctx context.Context, transferRequest transactionWeb.TransferRequest) (transactionWeb.TransactionResponse, error)

	FindAll(ctx context.Context, filter map[string]interface{}, paginationRequest httpservice.PaginationRequest, orderRequest httpservice.OrderRequest) ([]transactionWeb.TransactionResponse, httpservice.PaginationResponse, error)
}
