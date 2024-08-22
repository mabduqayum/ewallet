package services

import (
	"context"

	"github.com/mabduqayum/ewallet/internal/models"
	"github.com/mabduqayum/ewallet/internal/repository"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo repository.TransactionRepository
}

func NewTransactionService(repo repository.TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(ctx context.Context, walletID uuid.UUID, transactionType models.TransactionType, amount float64, description string) (*models.Transaction, error) {
	transaction := models.NewTransaction(walletID, transactionType, amount, description)
	err := s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (s *TransactionService) GetTransactionByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TransactionService) GetTransactionsByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	return s.repo.GetByWalletID(ctx, walletID, limit, offset)
}

func (s *TransactionService) GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error) {
	return s.repo.GetMonthlyTopUpStats(ctx, walletID)
}
