package repository

import (
	"context"
	"database/sql"
	"errors"
	"ewallet/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type WalletRepository interface {
	Exists(ctx context.Context, walletID uuid.UUID) (bool, error)
	GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error)
}

type PostgresWalletRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresWalletRepository(db *pgxpool.Pool) *PostgresWalletRepository {
	return &PostgresWalletRepository{pool: db}
}

func (r *PostgresWalletRepository) Exists(ctx context.Context, walletID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1)", walletID).Scan(&exists)
	return exists, err
}

func (r *PostgresWalletRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	err := r.pool.QueryRow(ctx, "SELECT id, type, balance, created_at, updated_at FROM wallets WHERE id = $1", walletID).
		Scan(&wallet.ID, &wallet.Type, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return wallet, err
}

func (r *PostgresWalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	_, err := r.pool.Exec(ctx, "UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3",
		wallet.Balance, time.Now(), wallet.ID)
	return err
}

func (r *PostgresWalletRepository) GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error) {
	var count int
	var sum float64
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE wallet_id = $1 AND type = 'top_up' AND created_at >= DATE_TRUNC('month', CURRENT_DATE)
	`, walletID).Scan(&count, &sum)
	return count, sum, err
}
