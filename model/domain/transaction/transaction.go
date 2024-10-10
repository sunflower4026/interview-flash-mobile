package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TOPUP    TransactionType = "TOPUP"
	PAYMENT  TransactionType = "PAYMENT"
	TRANSFER TransactionType = "TRANSFER"
)

type AccountType string

const (
	DEBIT  AccountType = "DEBIT"
	CREDIT AccountType = "CREDIT"
)

type Transaction struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;"`
	UserID          uuid.UUID `gorm:"type:uuid;"`
	TransactionType TransactionType
	AccountType     AccountType `gorm:"->"`
	Amount          int64
	Remarks         string
	CreatedAt       time.Time

	// Fields below are not stored in the database (view fields)
	BalanceBefore int64 `gorm:"->"`
	BalanceAfter  int64 `gorm:"->"`
}

func (u Transaction) TableName() string {
	return "transactions"
}

func (u Transaction) ViewTableName() string {
	return "view_transactions"
}

func (u *Transaction) BeforeCreate(*gorm.DB) error {
	u.ID = uuid.New()
	return nil
}
