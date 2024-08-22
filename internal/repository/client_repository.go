package repository

import (
	"context"

	"github.com/mabduqayum/ewallet/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClientRepository interface {
	Create(ctx context.Context, client *models.Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error)
	GetByAPIKey(ctx context.Context, apiKey string) (*models.Client, error)
	GetAll(ctx context.Context) ([]*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type PostgresClientRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresClientRepository(pool *pgxpool.Pool) *PostgresClientRepository {
	return &PostgresClientRepository{pool: pool}
}

func (r *PostgresClientRepository) Create(ctx context.Context, client *models.Client) error {
	_, err := r.pool.Exec(ctx,
		"INSERT INTO clients (id, name, api_key, secret_key, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		client.ID, client.Name, client.ApiKey, client.SecretKey, client.Active, client.CreatedAt, client.UpdatedAt)
	return err
}

func (r *PostgresClientRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	client := &models.Client{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, api_key, secret_key, active, created_at, updated_at FROM clients WHERE id = $1",
		id).Scan(&client.ID, &client.Name, &client.ApiKey, &client.SecretKey, &client.Active, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *PostgresClientRepository) GetByAPIKey(ctx context.Context, apiKey string) (*models.Client, error) {
	client := &models.Client{}
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, api_key, secret_key, active, created_at, updated_at FROM clients WHERE api_key = $1",
		apiKey).Scan(&client.ID, &client.Name, &client.ApiKey, &client.SecretKey, &client.Active, &client.CreatedAt, &client.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (r *PostgresClientRepository) GetAll(ctx context.Context) ([]*models.Client, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name, api_key, secret_key, active, created_at, updated_at FROM clients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []*models.Client
	for rows.Next() {
		client := &models.Client{}
		err := rows.Scan(&client.ID, &client.Name, &client.ApiKey, &client.SecretKey, &client.Active, &client.CreatedAt, &client.UpdatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func (r *PostgresClientRepository) Update(ctx context.Context, client *models.Client) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE clients SET name = $1, api_key = $2, secret_key = $3, active = $4, updated_at = $5 WHERE id = $6",
		client.Name, client.ApiKey, client.SecretKey, client.Active, client.UpdatedAt, client.ID)
	return err
}

func (r *PostgresClientRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM clients WHERE id = $1", id)
	return err
}
