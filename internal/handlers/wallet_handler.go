// internal/handlers/wallet_handler.go

package handlers

import (
	"ewallet/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}

func (h *WalletHandler) CheckWalletExists(c *fiber.Ctx) error {
	var req struct {
		WalletID string `json:"walletID"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid wallet ID"})
	}

	exists, err := h.walletService.CheckWalletExists(c.Context(), walletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check wallet existence"})
	}

	return c.JSON(fiber.Map{"exists": exists})
}

func (h *WalletHandler) TopUpWallet(c *fiber.Ctx) error {
	var req struct {
		WalletID string  `json:"walletID"`
		Amount   float64 `json:"amount"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid wallet ID"})
	}

	if req.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Amount must be positive"})
	}

	err = h.walletService.TopUpWallet(c.Context(), walletID, req.Amount)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Wallet topped up successfully"})
}

func (h *WalletHandler) GetMonthlyTopUpStats(c *fiber.Ctx) error {
	var req struct {
		WalletID string `json:"walletID"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid wallet ID"})
	}

	count, sum, err := h.walletService.GetMonthlyTopUpStats(c.Context(), walletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get monthly top-up stats"})
	}

	return c.JSON(fiber.Map{
		"count": count,
		"sum":   sum,
	})
}

func (h *WalletHandler) GetBalance(c *fiber.Ctx) error {
	var req struct {
		WalletID string `json:"walletID"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	walletID, err := uuid.Parse(req.WalletID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid wallet ID"})
	}

	balance, err := h.walletService.GetBalance(c.Context(), walletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get wallet balance"})
	}

	return c.JSON(fiber.Map{"balance": balance})
}
