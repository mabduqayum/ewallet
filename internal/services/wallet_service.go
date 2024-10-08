package services

import (
	"context"
	"errors"

	"github.com/mabduqayum/ewallet/internal/repository"

	"github.com/google/uuid"
)

type WalletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) CheckWalletExists(ctx context.Context, walletID uuid.UUID) (bool, error) {
	return s.repo.Exists(ctx, walletID)
}

func (s *WalletService) TopUpWallet(ctx context.Context, walletID uuid.UUID, amount float64) error {
	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return err
	}

	if wallet == nil {
		return errors.New("wallet not found")
	}

	err = wallet.UpdateBalance(amount)
	if err != nil {
		return err
	}

	return s.repo.Update(ctx, wallet, amount)
}

func (s *WalletService) GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error) {
	return s.repo.GetMonthlyTopUpStats(ctx, walletID)
}

func (s *WalletService) GetBalance(ctx context.Context, walletID uuid.UUID) (float64, error) {
	wallet, err := s.repo.GetByID(ctx, walletID)
	if err != nil {
		return 0, err
	}

	if wallet == nil {
		return 0, errors.New("wallet not found")
	}

	return wallet.Balance, nil
}
