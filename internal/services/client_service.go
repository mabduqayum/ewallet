package services

import (
	"context"

	"github.com/mabduqayum/ewallet/internal/models"
	"github.com/mabduqayum/ewallet/internal/repository"

	"github.com/google/uuid"
)

type ClientService struct {
	repo repository.ClientRepository
}

func NewClientService(repo repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

func (s *ClientService) CreateClient(ctx context.Context, name string) (*models.Client, error) {
	client := models.NewClient(name)
	err := s.repo.Create(ctx, client)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) GetClientByID(ctx context.Context, id uuid.UUID) (*models.Client, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ClientService) GetClientByAPIKey(ctx context.Context, apiKey string) (*models.Client, error) {
	return s.repo.GetByAPIKey(ctx, apiKey)
}

func (s *ClientService) GetAllClients(ctx context.Context) ([]*models.Client, error) {
	return s.repo.GetAll(ctx)
}

func (s *ClientService) UpdateClient(ctx context.Context, client *models.Client) error {
	return s.repo.Update(ctx, client)
}

func (s *ClientService) DeleteClient(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
