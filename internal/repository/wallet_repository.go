package repository

import (
	"context"
	"database/sql"
	"errors"
	"ewallet/internal/models"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet models.Wallet) error
	Exists(ctx context.Context, walletID uuid.UUID) (bool, error)
	GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet) error
	GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error)
}

type PostgresWalletRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresWalletRepository(pool *pgxpool.Pool) *PostgresWalletRepository {
	return &PostgresWalletRepository{pool: pool}
}

func (r *PostgresWalletRepository) Create(ctx context.Context, wallet models.Wallet) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO wallets (id, type, balance, currency, active, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		wallet.ID,
		wallet.Type,
		wallet.Balance,
		wallet.Currency,
		wallet.Active,
		wallet.CreatedAt,
		wallet.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	return nil
}

func (r *PostgresWalletRepository) Exists(ctx context.Context, walletID uuid.UUID) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1)", walletID).Scan(&exists)
	return exists, err
}

func (r *PostgresWalletRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error) {
	wallet := &models.Wallet{}
	err := r.pool.QueryRow(ctx, "SELECT id, type, balance, currency, active, created_at, updated_at FROM wallets WHERE id = $1", walletID).
		Scan(&wallet.ID, &wallet.Type, &wallet.Balance, &wallet.Currency, &wallet.Active, &wallet.CreatedAt, &wallet.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return wallet, err
}

func (r *PostgresWalletRepository) Update(ctx context.Context, wallet *models.Wallet) error {
	_, err := r.pool.Exec(ctx, "UPDATE wallets SET balance = $1, active = $2, updated_at = $3 WHERE id = $4",
		wallet.Balance, wallet.Active, time.Now(), wallet.ID)
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
