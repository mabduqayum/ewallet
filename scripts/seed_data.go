package scripts

import (
	"context"
	"ewallet/internal/models"
	"ewallet/internal/repository"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func SeedData(ctx context.Context, pool *pgxpool.Pool) error {
	clientRepo := repository.NewPostgresClientRepository(pool)
	walletRepo := repository.NewPostgresWalletRepository(pool)
	transactionRepo := repository.NewPostgresTransactionRepository(pool)

	clients, err := seedClients(ctx, clientRepo)
	if err != nil {
		return err
	}

	wallets, err := seedWallets(ctx, walletRepo, clients)
	if err != nil {
		return err
	}

	err = seedTransactions(ctx, transactionRepo, walletRepo, wallets)
	if err != nil {
		return err
	}

	return nil
}

func seedClients(ctx context.Context, repo repository.ClientRepository) ([]*models.Client, error) {
	clients := []*models.Client{
		models.NewClient("Client 1"),
		models.NewClient("Client 2"),
		models.NewClient("Client 3"),
	}

	for _, client := range clients {
		err := repo.Create(ctx, client)
		if err != nil {
			return nil, err
		}
	}

	log.Printf("Seeded %d clients", len(clients))
	return clients, nil
}

func seedWallets(ctx context.Context, repo repository.WalletRepository, clients []*models.Client) ([]*models.Wallet, error) {
	wallets := make([]*models.Wallet, 0)

	for _, client := range clients {
		identifiedWallet := models.NewWallet(client.ID, models.WalletTypeIdentified, "TJS")
		unidentifiedWallet := models.NewWallet(client.ID, models.WalletTypeUnidentified, "TJS")

		wallets = append(wallets, identifiedWallet, unidentifiedWallet)

		for _, wallet := range []*models.Wallet{identifiedWallet, unidentifiedWallet} {
			err := repo.Create(ctx, *wallet)
			if err != nil {
				return nil, err
			}
		}
	}

	log.Printf("Seeded %d wallets", len(wallets))
	return wallets, nil
}

func seedTransactions(ctx context.Context, transactionRepo repository.TransactionRepository, walletRepo repository.WalletRepository, wallets []*models.Wallet) error {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	for _, wallet := range wallets {
		numTransactions := r.Intn(10) + 1 // 1 to 10 transactions per wallet

		for i := 0; i < numTransactions; i++ {
			amount := r.Float64() * 1000 // Random amount up to 1000
			transaction := models.NewTransaction(wallet.ID, models.TransactionTypeTopUp, amount, "Initial top-up")

			err := transactionRepo.Create(ctx, transaction)
			if err != nil {
				return err
			}

			// Update wallet balance
			wallet.Balance += amount
			err = walletRepo.Update(ctx, wallet)
			if err != nil {
				return err
			}
		}
	}

	log.Println("Seeded transactions")
	return nil
}
