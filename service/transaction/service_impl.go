package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/constants"
	"gitlab.com/sunflower4026/interview-flash-mobile/common/httpservice"
	transactionDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/transaction"
	transactionRepository "gitlab.com/sunflower4026/interview-flash-mobile/model/repository/transaction"
	transactionWeb "gitlab.com/sunflower4026/interview-flash-mobile/model/web/transaction"
	"gitlab.com/sunflower4026/interview-flash-mobile/toolkit/log"
	"gorm.io/gorm"
)

type TransactionServiceImpl struct {
	TransactionRepository transactionRepository.TransactionRepository
	DB                    *gorm.DB
}

func NewTransactionService(transactionRepository transactionRepository.TransactionRepository, DB *gorm.DB) TransactionService {
	return &TransactionServiceImpl{
		TransactionRepository: transactionRepository,
		DB:                    DB,
	}
}

func (s *TransactionServiceImpl) Topup(ctx context.Context, topupRequest transactionWeb.TopupRequest) (transactionWeb.TransactionResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		log.FromCtx(ctx).Error(errors.New("user id not found in context"), "user id not found in context")
		return transactionWeb.TransactionResponse{}, errors.New("user id not found in context")
	}

	transaction := topupRequest.ToTransaction()
	transaction.UserID = uuid.MustParse(userID)
	transaction.TransactionType = transactionDomain.TOPUP

	err := s.TransactionRepository.Create(ctx, tx, &transaction)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	transaction, err = s.TransactionRepository.Find(ctx, tx, fmt.Sprintf("id = '%s'", transaction.ID.String()))
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	tx.Commit()

	return transactionWeb.ToTransactionResponse(transaction), nil
}

func (s *TransactionServiceImpl) Payment(ctx context.Context, paymentRequest transactionWeb.PaymentRequest) (transactionWeb.TransactionResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		log.FromCtx(ctx).Error(errors.New("user id not found in context"), "user id not found in context")
		return transactionWeb.TransactionResponse{}, errors.New("user id not found in context")
	}

	// Check if user has enough balance
	userBalance, _, err := s.TransactionRepository.FindAll(ctx, tx, map[string]interface{}{"user_id": userID}, httpservice.PaginationRequest{Page: 1, Limit: 1}, httpservice.OrderRequest{OrderBy: "created_at", Sort: "desc"})
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find user balance")
		return transactionWeb.TransactionResponse{}, err
	}

	if len(userBalance) == 0 {
		log.FromCtx(ctx).Error(errors.New("user balance not found"), "user balance not found")
		return transactionWeb.TransactionResponse{}, errors.New("user balance not found")
	}

	if userBalance[0].BalanceAfter < int64(paymentRequest.Amount) {
		log.FromCtx(ctx).Error(errors.New("user has not enough balance"), "user has not enough balance")
		return transactionWeb.TransactionResponse{}, errors.New("user has not enough balance")
	}

	transaction := paymentRequest.ToTransaction()
	transaction.UserID = uuid.MustParse(userID)
	transaction.TransactionType = transactionDomain.PAYMENT

	err = s.TransactionRepository.Create(ctx, tx, &transaction)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	transaction, err = s.TransactionRepository.Find(ctx, tx, fmt.Sprintf("id = '%s'", transaction.ID.String()))
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	tx.Commit()

	return transactionWeb.ToTransactionResponse(transaction), nil
}

func (s *TransactionServiceImpl) Transfer(ctx context.Context, transferRequest transactionWeb.TransferRequest) (transactionWeb.TransactionResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		log.FromCtx(ctx).Error(errors.New("user id not found in context"), "user id not found in context")
		return transactionWeb.TransactionResponse{}, errors.New("user id not found in context")
	}

	// Check if user has enough balance
	userBalance, _, err := s.TransactionRepository.FindAll(ctx, tx, map[string]interface{}{"user_id": userID}, httpservice.PaginationRequest{Page: 1, Limit: 1}, httpservice.OrderRequest{OrderBy: "created_at", Sort: "desc"})
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find user balance")
		return transactionWeb.TransactionResponse{}, err
	}

	if len(userBalance) == 0 {
		log.FromCtx(ctx).Error(errors.New("user balance not found"), "user balance not found")
		return transactionWeb.TransactionResponse{}, errors.New("user balance not found")
	}

	if userBalance[0].BalanceAfter < int64(transferRequest.Amount) {
		log.FromCtx(ctx).Error(errors.New("user has not enough balance"), "user has not enough balance")
		return transactionWeb.TransactionResponse{}, errors.New("user has not enough balance")
	}

	transaction := transferRequest.ToTransactionSender()
	transaction.UserID = uuid.MustParse(userID)
	transaction.TransactionType = transactionDomain.TRANSFER

	err = s.TransactionRepository.Create(ctx, tx, &transaction)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	transactionReceiver := transferRequest.ToTransactionReceiver()
	transactionReceiver.UserID = uuid.MustParse(transferRequest.To)
	transactionReceiver.TransactionType = transactionDomain.TRANSFER

	err = s.TransactionRepository.Create(ctx, tx, &transactionReceiver)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to create transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	transaction, err = s.TransactionRepository.Find(ctx, tx, fmt.Sprintf("id = '%s'", transaction.ID.String()))
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find transaction")
		return transactionWeb.TransactionResponse{}, err
	}

	tx.Commit()

	return transactionWeb.ToTransactionResponse(transaction), nil
}

func (s *TransactionServiceImpl) FindAll(ctx context.Context, filter map[string]interface{}, paginationRequest httpservice.PaginationRequest, orderRequest httpservice.OrderRequest) ([]transactionWeb.TransactionResponse, httpservice.PaginationResponse, error) {
	tx := s.DB.Begin()
	defer tx.Rollback()

	userID, ok := ctx.Value(constants.ContextKeyUserID).(string)
	if !ok {
		log.FromCtx(ctx).Error(errors.New("user id not found in context"), "user id not found in context")
		return nil, httpservice.PaginationResponse{}, errors.New("user id not found in context")
	}

	filter["user_id"] = userID

	transactions, paginationResponse, err := s.TransactionRepository.FindAll(ctx, tx, filter, paginationRequest, orderRequest)
	if err != nil {
		log.FromCtx(ctx).Error(err, "Failed to find all transactions")
		return nil, httpservice.PaginationResponse{}, err
	}

	tx.Commit()

	return transactionWeb.ToTransactionResponses(transactions), paginationResponse, nil
}
