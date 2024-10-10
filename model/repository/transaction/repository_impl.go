package transaction

import (
	"context"

	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"

	transactionDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/transaction"
	"gorm.io/gorm"
)

type TransactionRepositoryImpl struct {
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImpl{}
}

func (r *TransactionRepositoryImpl) Create(ctx context.Context, tx *gorm.DB, transaction *transactionDomain.Transaction) error {
	if err := tx.Create(transaction).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create transaction")
		return err
	}

	return nil
}

func (r *TransactionRepositoryImpl) Find(ctx context.Context, tx *gorm.DB, where string) (transactionDomain.Transaction, error) {
	var transaction transactionDomain.Transaction
	if err := tx.Table(transactionDomain.Transaction{}.ViewTableName()).Where(where).First(&transaction).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find transaction")
		return transaction, err
	}

	return transaction, nil
}

func (r *TransactionRepositoryImpl) FindAll(ctx context.Context, tx *gorm.DB, filter map[string]interface{}, paginationRequest httpservice.PaginationRequest, orderRequest httpservice.OrderRequest) ([]transactionDomain.Transaction, httpservice.PaginationResponse, error) {
	var transactions []transactionDomain.Transaction
	var count int64

	tx = tx.Table(transactionDomain.Transaction{}.ViewTableName())

	for key, value := range filter {
		tx = tx.Where(key, value)
	}

	if err := tx.Count(&count).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to count transactions")
		return nil, httpservice.PaginationResponse{}, err
	}

	if err := tx.Limit(paginationRequest.Limit).Offset((paginationRequest.Page - 1) * paginationRequest.Limit).Order(orderRequest.OrderBy + " " + orderRequest.Sort).Find(&transactions).Error; err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find transactions")
		return nil, httpservice.PaginationResponse{}, err
	}

	return transactions, httpservice.PaginationResponse{
		TotalPage: int(count) / paginationRequest.Limit,
		TotalData: count,
	}, nil
}
