package types

import "github.com/gofiber/fiber/v3"

type Controller interface {
	RegisterEndpoint(router fiber.Router)
}
