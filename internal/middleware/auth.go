package middleware

import (
	"ewallet/internal/services"
	"ewallet/internal/utils/hmac"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(clientService *services.ClientService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Get("X-UserId")
		digest := c.Get("X-Digest")

		if userID == "" || digest == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authentication headers",
			})
		}

		client, err := clientService.GetClientByAPIKey(c.Context(), userID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID",
			})
		}

		body := c.Body()
		if !hmac.ValidateHMAC(string(body), client.SecretKey, digest) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid digest",
			})
		}

		c.Request().SetBody(body)

		return c.Next()
	}
}
