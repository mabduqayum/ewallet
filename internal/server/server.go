package server

import (
	"github.com/gofiber/fiber/v2"

	"ewallet/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "ewallet",
			AppName:      "ewallet",
		}),

		db: database.New(),
	}

	return server
}
