package transaction

import (
	"context"

	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	transactionDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/transaction"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *gorm.DB, transaction *transactionDomain.Transaction) error
	Find(ctx context.Context, tx *gorm.DB, whr string) (transactionDomain.Transaction, error)
	FindAll(ctx context.Context, tx *gorm.DB, filter map[string]interface{}, paginationRequest httpservice.PaginationRequest, orderRequest httpservice.OrderRequest) ([]transactionDomain.Transaction, httpservice.PaginationResponse, error)
}
