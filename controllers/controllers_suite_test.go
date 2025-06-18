package controllers_test

import (
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/hotrungnhan/surl/utils/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

func AppTest() (*fiber.App, fiber.Router) {
	app := fiber.New(
		fiber.Config{
			ErrorHandler:             types.CustomErrorHandler,
			EnableSplittingOnParsers: true,
		},
	)

	router := app.Group("")
	return app, router
}
