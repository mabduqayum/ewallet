package server

import (
	"ewallet/internal/config"
	"ewallet/internal/database"
	"ewallet/internal/repository"
	"ewallet/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type FiberServer struct {
	app *fiber.App
	db  database.Service
	cfg *config.Config

	walletService *services.WalletService
}

func New(cfg *config.Config, db database.Service) *FiberServer {
	walletRepo := repository.NewPostgresWalletRepository(db.GetPool())
	walletService := services.NewWalletService(walletRepo)

	server := &FiberServer{
		app: fiber.New(fiber.Config{
			ServerHeader: "ewallet",
			AppName:      "ewallet v" + cfg.Server.Version,
		}),

		db:            db,
		cfg:           cfg,
		walletService: walletService,
	}

	// Add recover middleware
	server.app.Use(recover.New())

	return server
}

func (s *FiberServer) Listen() error {
	return s.app.Listen(s.cfg.Server.Address())
}

func (s *FiberServer) Shutdown() error {
	return s.app.Shutdown()
}
