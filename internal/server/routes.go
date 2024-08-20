package server

import (
	"context"
	"ewallet/internal/handlers"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/contrib/websocket"
)

func (s *FiberServer) RegisterFiberRoutes() {
	s.app.Get("/", s.HelloWorldHandler)
	s.app.Get("/health", s.healthHandler)
	s.app.Get("/websocket", websocket.New(s.websocketHandler))

	api := s.app.Group("/api/v1")

	wallet := api.Group("/wallet")
	walletHandler := handlers.NewWalletHandler(s.walletService)
	wallet.Post("/exists", walletHandler.CheckWalletExists)
	wallet.Post("/top-up", walletHandler.TopUpWallet)
	wallet.Post("/stats", walletHandler.GetMonthlyTopUpStats)
	wallet.Post("/balance", walletHandler.GetBalance)
}

func (s *FiberServer) HelloWorldHandler(c *fiber.Ctx) error {
	resp := fiber.Map{
		"message": "Hello World",
	}

	return c.JSON(resp)
}

func (s *FiberServer) healthHandler(c *fiber.Ctx) error {
	return c.JSON(s.db.Health())
}

func (s *FiberServer) websocketHandler(con *websocket.Conn) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			_, _, err := con.ReadMessage()
			if err != nil {
				cancel()
				log.Println("Receiver Closing", err)
				break
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
			if err := con.WriteMessage(websocket.TextMessage, []byte(payload)); err != nil {
				log.Printf("could not write to socket: %v", err)
				return
			}
			time.Sleep(time.Second * 2)
		}
	}
}
