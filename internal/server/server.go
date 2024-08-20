package server

import (
	"ewallet/internal/config"
	"ewallet/internal/database"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App

	db  database.Service
	cfg *config.Config
}

func New(cfg *config.Config, db database.Service) *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ewallet",
			AppName:      "ewallet v" + cfg.Server.Version,
		}),

		db:  db,
		cfg: cfg,
	}

	return server
}

func (s *FiberServer) Listen() error {
	return s.App.Listen(s.cfg.Server.Address())
}

func (s *FiberServer) Shutdown() error {
	return s.App.Shutdown()
}
