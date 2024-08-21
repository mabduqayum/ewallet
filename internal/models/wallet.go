package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WalletType string

const (
	WalletTypeIdentified   WalletType = "IDENTIFIED"
	WalletTypeUnidentified WalletType = "UNIDENTIFIED"
)

type Wallet struct {
	ID        uuid.UUID  `json:"id"`
	Type      WalletType `json:"type"`
	Balance   float64    `json:"balance"`
	Currency  string     `json:"currency"`
	Active    bool       `json:"active"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewWallet(walletType WalletType, currency string) *Wallet {
	return &Wallet{
		ID:        uuid.New(),
		Type:      walletType,
		Balance:   0,
		Currency:  currency,
		Active:    false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (w *Wallet) GetMaxBalance() float64 {
	switch w.Type {
	case WalletTypeIdentified:
		return 100_000
	case WalletTypeUnidentified:
		return 10_000
	default:
		return 0
	}
}

func (w *Wallet) UpdateBalance(amount float64) error {
	newBalance := w.Balance + amount
	if newBalance < 0 {
		return errors.New("insufficient funds")
	}
	if newBalance > w.GetMaxBalance() {
		return errors.New("balance exceeds maximum limit")
	}
	w.Balance = newBalance
	w.UpdatedAt = time.Now()
	return nil
}

func (w *Wallet) Activate() {
	w.Active = true
	w.UpdatedAt = time.Now()
}
