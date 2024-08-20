package scripts

import (
	"context"
	"ewallet/internal/database"
	"ewallet/internal/models"
	"log"
	"math/rand"
	"time"
)

func SeedClients(ctx context.Context, db database.Service) []*models.Client {
	clients := []*models.Client{
		models.NewClient("Client 1"),
		models.NewClient("Client 2"),
		models.NewClient("Client 3"),
	}

	for _, client := range clients {
		_, err := db.GetPool().Exec(ctx,
			"INSERT INTO clients (id, name, api_key, secret_key, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			client.ID, client.Name, client.ApiKey, client.SecretKey, client.Active, client.CreatedAt, client.UpdatedAt)
		if err != nil {
			log.Fatalf("Failed to insert client: %v", err)
		}
	}

	log.Printf("Seeded %d clients", len(clients))
	return clients
}

func SeedWallets(ctx context.Context, db database.Service, clients []*models.Client) []*models.Wallet {
	wallets := make([]*models.Wallet, 0)

	for _, client := range clients {
		identifiedWallet := models.NewWallet(client.ID, models.WalletTypeIdentified, "TJS")
		unidentifiedWallet := models.NewWallet(client.ID, models.WalletTypeUnidentified, "TJS")

		wallets = append(wallets, identifiedWallet, unidentifiedWallet)

		for _, wallet := range []*models.Wallet{identifiedWallet, unidentifiedWallet} {
			_, err := db.GetPool().Exec(ctx,
				"INSERT INTO wallets (id, client_id, type, balance, currency, active, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
				wallet.ID, wallet.ClientID, wallet.Type, wallet.Balance, wallet.Currency, wallet.Active, wallet.CreatedAt, wallet.UpdatedAt)
			if err != nil {
				log.Fatalf("Failed to insert wallet: %v", err)
			}
		}
	}

	log.Printf("Seeded %d wallets", len(wallets))
	return wallets
}

func SeedTransactions(ctx context.Context, db database.Service, wallets []*models.Wallet) {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for _, wallet := range wallets {
		numTransactions := r.Intn(10) + 1 // 1 to 10 transactions per wallet

		for i := 0; i < numTransactions; i++ {
			amount := r.Float64() * 1000 // Random amount up to 1000
			transaction := models.NewTransaction(wallet.ID, models.TransactionTypeTopUp, amount, "Initial top-up")

			_, err := db.GetPool().Exec(ctx,
				"INSERT INTO transactions (id, wallet_id, type, amount, description, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
				transaction.ID, transaction.WalletID, transaction.Type, transaction.Amount, transaction.Description, transaction.CreatedAt, transaction.UpdatedAt)
			if err != nil {
				log.Fatalf("Failed to insert transaction: %v", err)
			}

			// Update wallet balance
			_, err = db.GetPool().Exec(ctx,
				"UPDATE wallets SET balance = balance + $1, updated_at = $2 WHERE id = $3",
				transaction.Amount, time.Now(), wallet.ID)
			if err != nil {
				log.Fatalf("Failed to update wallet balance: %v", err)
			}
		}
	}

	log.Println("Seeded transactions")
}
