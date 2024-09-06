package server

import (
	"github.com/mabduqayum/ewallet/internal/config"
	"github.com/mabduqayum/ewallet/internal/database"
	"github.com/mabduqayum/ewallet/internal/repository"
	"github.com/mabduqayum/ewallet/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type FiberServer struct {
	app *fiber.App
	db  database.Service
	cfg *config.ServerConfig

	walletService *services.WalletService
	clientService *services.ClientService
}

func New(cfg *config.ServerConfig, db database.Service) *FiberServer {
	walletRepo := repository.NewPostgresWalletRepository(db.GetPool())
	walletService := services.NewWalletService(walletRepo)

	clientRepo := repository.NewPostgresClientRepository(db.GetPool())
	clientService := services.NewClientService(clientRepo)

	server := &FiberServer{
		app: fiber.New(fiber.Config{
			ServerHeader: "ewallet",
			AppName:      "ewallet v" + cfg.Version,
		}),

		db:            db,
		cfg:           cfg,
		walletService: walletService,
		clientService: clientService,
	}

	// Add recover middleware
	server.app.Use(recover.New())

	return server
}

func (s *FiberServer) Listen() error {
	return s.app.Listen(s.cfg.Address())
}

func (s *FiberServer) Shutdown() error {
	return s.app.Shutdown()
}
