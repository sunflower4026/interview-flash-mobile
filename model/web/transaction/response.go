package transaction

import (
	"github.com/google/uuid"

	transactionDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/transaction"
)

type TransactionResponse struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TransactionType string    `json:"transaction_type"`
	AccountType     string    `json:"account_type"`
	Amount          int64     `json:"amount"`
	Remarks         string    `json:"remarks"`
	CreatedAt       string    `json:"created_at"`

	BalanceBefore int64 `json:"balance_before"`
	BalanceAfter  int64 `json:"balance_after"`
}

func ToTransactionResponse(transaction transactionDomain.Transaction) TransactionResponse {
	return TransactionResponse{
		ID:              transaction.ID,
		UserID:          transaction.UserID,
		TransactionType: string(transaction.TransactionType),
		AccountType:     string(transaction.AccountType),
		Amount:          transaction.Amount,
		Remarks:         transaction.Remarks,
		CreatedAt:       transaction.CreatedAt.Format("2006-01-02 15:04:05"),

		BalanceBefore: transaction.BalanceBefore,
		BalanceAfter:  transaction.BalanceAfter,
	}
}

func ToTransactionResponses(transactions []transactionDomain.Transaction) []TransactionResponse {
	var transactionResponses []TransactionResponse
	for _, transaction := range transactions {
		transactionResponses = append(transactionResponses, ToTransactionResponse(transaction))
	}

	return transactionResponses
}
