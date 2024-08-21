package scripts

import (
	"context"
	"ewallet/internal/models"
	"ewallet/internal/repository"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	numClients     = 10
	numWallets     = 30
	maxTopUps      = 10
	maxTopUpAmount = 1000.0
)

func SeedData(ctx context.Context, pool *pgxpool.Pool) error {
	clientRepo := repository.NewPostgresClientRepository(pool)
	walletRepo := repository.NewPostgresWalletRepository(pool)
	transactionRepo := repository.NewPostgresTransactionRepository(pool)

	clients, err := seedClients(ctx, clientRepo)
	if err != nil {
		return err
	}

	wallets, err := seedWallets(ctx, walletRepo)
	if err != nil {
		return err
	}

	err = seedTransactions(ctx, transactionRepo, walletRepo, wallets)
	if err != nil {
		return err
	}

	log.Printf("Seeded %d clients, %d wallets, and transactions", len(clients), len(wallets))
	return nil
}

func seedClients(ctx context.Context, repo repository.ClientRepository) ([]*models.Client, error) {
	clients := make([]*models.Client, 0, numClients)

	for i := 1; i <= numClients; i++ {
		client := models.NewClient(fmt.Sprintf("Client %d", i))
		if err := repo.Create(ctx, client); err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func seedWallets(ctx context.Context, repo repository.WalletRepository) ([]*models.Wallet, error) {
	wallets := make([]*models.Wallet, 0, numWallets)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range numWallets {
		walletType := models.WalletTypeIdentified
		if r.Float32() < 0.5 {
			walletType = models.WalletTypeUnidentified
		}

		wallet := models.NewWallet(walletType, "TJS")
		if err := repo.Create(ctx, *wallet); err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}

	return wallets, nil
}

func seedTransactions(ctx context.Context, transactionRepo repository.TransactionRepository, walletRepo repository.WalletRepository, wallets []*models.Wallet) error {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, wallet := range wallets {
		numTransactions := r.Intn(maxTopUps) + 1

		for i := 0; i < numTransactions; i++ {
			amount := r.Float64() * maxTopUpAmount
			transaction := models.NewTransaction(wallet.ID, models.TransactionTypeTopUp, amount, "Initial top-up")

			if err := transactionRepo.Create(ctx, transaction); err != nil {
				return err
			}

			wallet.Balance += amount
			if err := walletRepo.Update(ctx, wallet, amount); err != nil {
				return err
			}
		}
	}

	return nil
}
