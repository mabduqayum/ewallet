package models

import (
	"time"

	"github.com/google/uuid"
)

type TransactionType string

const (
	TransactionTypeTopUp TransactionType = "TOP_UP"
	// TransactionTypeWithdraw TransactionType = "WITHDRAW"
)

type Transaction struct {
	ID          uuid.UUID       `json:"id"`
	WalletID    uuid.UUID       `json:"wallet_id"`
	Type        TransactionType `json:"type"`
	Amount      float64         `json:"amount"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func NewTransaction(walletID uuid.UUID, transactionType TransactionType, amount float64, description string) *Transaction {
	return &Transaction{
		ID:          uuid.New(),
		WalletID:    walletID,
		Type:        transactionType,
		Amount:      amount,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
