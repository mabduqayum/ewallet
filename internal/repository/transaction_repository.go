package repository

import (
	"context"
	"ewallet/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error)
	GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Transaction, error)
	GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error)
}

type PostgresTransactionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresTransactionRepository(pool *pgxpool.Pool) *PostgresTransactionRepository {
	return &PostgresTransactionRepository{pool: pool}
}

func (r *PostgresTransactionRepository) Create(ctx context.Context, transaction *models.Transaction) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO transactions (id, wallet_id, type, amount, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		transaction.ID, transaction.WalletID, transaction.Type, transaction.Amount, transaction.Description, transaction.CreatedAt, transaction.UpdatedAt)
	return err
}

func (r *PostgresTransactionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Transaction, error) {
	transaction := &models.Transaction{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, wallet_id, type, amount, description, created_at, updated_at FROM transactions WHERE id = $1",
		id).Scan(&transaction.ID, &transaction.WalletID, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *PostgresTransactionRepository) GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*models.Transaction, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, wallet_id, type, amount, description, created_at, updated_at FROM transactions WHERE wallet_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3",
		walletID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*models.Transaction
	for rows.Next() {
		transaction := &models.Transaction{}
		err := rows.Scan(&transaction.ID, &transaction.WalletID, &transaction.Type, &transaction.Amount, &transaction.Description, &transaction.CreatedAt, &transaction.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *PostgresTransactionRepository) GetMonthlyTopUpStats(ctx context.Context, walletID uuid.UUID) (int, float64, error) {
	var count int
	var sum float64
	err := r.pool.QueryRow(ctx,
		"SELECT COUNT(*), COALESCE(SUM(amount), 0) FROM transactions WHERE wallet_id = $1 AND type = $2 AND created_at >= DATE_TRUNC('month', CURRENT_DATE)",
		walletID, models.TransactionTypeTopUp).Scan(&count, &sum)
	return count, sum, err
}
