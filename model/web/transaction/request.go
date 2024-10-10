package transaction

import (
	transactionDomain "gitlab.com/sunflower4026/interview-flash-mobile/model/domain/transaction"
)

type TopupRequest struct {
	Amount int `json:"amount"`
}

func (t TopupRequest) ToTransaction() transactionDomain.Transaction {
	return transactionDomain.Transaction{
		Amount: int64(t.Amount),
	}
}

type PaymentRequest struct {
	Amount  int    `json:"amount"`
	Remarks string `json:"remarks"`
}

func (p PaymentRequest) ToTransaction() transactionDomain.Transaction {
	return transactionDomain.Transaction{
		Amount:  int64(-p.Amount),
		Remarks: p.Remarks,
	}
}

type TransferRequest struct {
	To      string `json:"to"`
	Amount  int    `json:"amount"`
	Remarks string `json:"remarks"`
}

func (t TransferRequest) ToTransactionSender() transactionDomain.Transaction {
	return transactionDomain.Transaction{
		Amount:  int64(-t.Amount),
		Remarks: t.Remarks,
	}
}

func (t TransferRequest) ToTransactionReceiver() transactionDomain.Transaction {
	return transactionDomain.Transaction{
		Amount:  int64(t.Amount),
		Remarks: t.Remarks,
	}
}
