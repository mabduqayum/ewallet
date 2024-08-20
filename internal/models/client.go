package models

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
	SecretKey string    `json:"secret_key"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewClient(name string) *Client {
	return &Client{
		ID:        uuid.New(),
		Name:      name,
		ApiKey:    uuid.New().String(),
		SecretKey: uuid.New().String(),
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
