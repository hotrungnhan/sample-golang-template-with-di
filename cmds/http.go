package cmds

import (
	"fmt"
	"github.com/hotrungnhan/surl/controllers"
	"github.com/hotrungnhan/surl/repositories"
	"github.com/hotrungnhan/surl/services"
	"github.com/hotrungnhan/surl/utils/injects"
	"github.com/hotrungnhan/surl/utils/types"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/compress"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/requestid"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

type HttpServerOptions struct {
	Address string
	Port    int
}

func (o *HttpServerOptions) GetFullAddress() string {
	return o.Address + ":" + fmt.Sprint(o.Port)
}

type HttpServer struct {
	Options HttpServerOptions
}

func NewHttpServer(opts HttpServerOptions) types.Service {
	return &HttpServer{Options: opts}
}

// Serve implements Service.
func (d *HttpServer) Serve() error {

	di := do.New()

	// Base
	do.Provide(di, injects.NewAppConfig)
	do.Provide(di, injects.NewCache)
	do.Provide(di, injects.NewLogger)
	do.Provide(di, injects.NewDatabase)

	// Repository
	do.Provide(di, repositories.NewTransactionRepository)
	do.Provide(di, repositories.NewShortenUrlRepository)

	// Services
	do.Provide(di, services.NewShortenUrlService)

	// Controllers
	do.Provide(di, controllers.NewShortenUrlController)

	logger := do.MustInvoke[*zerolog.Logger](di)
	config := do.MustInvoke[injects.ApplicationConfig](di)

	logger.Info().Msgf("Starting HTTP server on %s", d.Options.GetFullAddress())

	app := fiber.New(fiber.Config{
		ErrorHandler:             types.CustomErrorHandler,
		EnableSplittingOnParsers: true,
	})

	if config.GoEnv == injects.Production {
		app.Use(requestid.New())

		app.Use(compress.New())

		app.Use(cors.New())
	}

	d.RegisterRouter(app, di)

	app.Hooks().OnShutdown(func() error {
		logger.Info().Msg("HTTP server is shutting down")
		return nil
	})

	return app.Listen(d.Options.GetFullAddress(), fiber.ListenConfig{
		EnablePrefork:     false,
		EnablePrintRoutes: true,
	})
}

func (d *HttpServer) RegisterRouter(app *fiber.App, di do.Injector) {

	base := app.Group("")
	base.Get("/live", healthcheck.NewHealthChecker()).Name("Live healcheck")

	{
		v1Controllers := []types.Controller{
			do.MustInvoke[controllers.ShortenUrlController](di),
		}

		lo.ForEach(v1Controllers, func(controller types.Controller, _ int) {
			controller.RegisterEndpoint(base)
		})
	}

}
func NewHttpServerCommand() *cli.Command {
	return &cli.Command{
		Name:  "http",
		Usage: "Run Http server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "address",
				Aliases: []string{"a"},
				Usage:   "Address to bind the DNS server to",
				Value:   "0.0.0.0",
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Usage:   "Port to run the DNS server on",
				Value:   8080,
			},
		},
		Action: func(c *cli.Context) error {
			address := c.String("address")
			port := c.Int("port")

			return NewHttpServer(HttpServerOptions{
				Address: address,
				Port:    port,
			}).Serve()
		},
	}
}
