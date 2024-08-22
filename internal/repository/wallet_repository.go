package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/mabduqayum/ewallet/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepository interface {
	Create(ctx context.Context, wallet models.Wallet) error
	Exists(ctx context.Context, walletID uuid.UUID) (bool, error)
	GetByID(ctx context.Context, walletID uuid.UUID) (*models.Wallet, error)
	Update(ctx context.Context, wallet *models.Wallet, amount float64) error
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
		`INSERT INTO wallets (id, type, balance, currency, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5, $6)`,
		wallet.ID,
		wallet.Type,
		wallet.Balance,
		wallet.Currency,
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
	err := r.pool.QueryRow(ctx, "SELECT id, type, balance, currency, created_at, updated_at FROM wallets WHERE id = $1", walletID).
		Scan(&wallet.ID, &wallet.Type, &wallet.Balance, &wallet.Currency, &wallet.CreatedAt, &wallet.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return wallet, err
}

func (r *PostgresWalletRepository) Update(ctx context.Context, wallet *models.Wallet, amount float64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		"UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3",
		wallet.Balance, time.Now(), wallet.ID)
	if err != nil {
		return fmt.Errorf("failed to update wallet balance: %w", err)
	}

	transaction := models.NewTransaction(wallet.ID, models.TransactionTypeTopUp, amount, "Top-up")
	_, err = tx.Exec(ctx, `
			INSERT INTO transactions (id, wallet_id, type, amount, description, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, transaction.ID, transaction.WalletID, transaction.Type, transaction.Amount, transaction.Description, transaction.CreatedAt, transaction.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create transaction record: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresWalletRepository) GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error) {
	var count int
	var sum float64
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)
		FROM transactions
		WHERE wallet_id = $1
		  AND type = 'TOP_UP'
		  AND created_at >= DATE_TRUNC('month', CURRENT_DATE)
	`, walletID).Scan(&count, &sum)
	return count, sum, err
}
