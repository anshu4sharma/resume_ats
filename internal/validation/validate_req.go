package validation

import (
	"github.com/gofiber/fiber/v2"
)

func ValidateBody[T any](validator *XValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body T
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid request body",
			})
		}

		errors := validator.Validate(&body)
		if len(errors) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":  true,
				"errors": errors,
			})
		}

		// Store validated body in context
		c.Locals("validated_body", body)
		return c.Next()
	}
}
