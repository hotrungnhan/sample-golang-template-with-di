package cmds

import (
	"fmt"
	"github.com/hotrungnhan/surl/utils/injects"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
	"strings"
)

func seedExample(tx *gorm.DB) error {
	return nil
}

func SeedHandler(di do.Injector, handlers map[string]func(tx *gorm.DB) error) error {
	db := do.MustInvoke[*injects.DB](di)
	logger := do.MustInvoke[*zerolog.Logger](di)

	logger.Info().Msg("Starting database seeding...")
	return db.Master.Transaction(func(tx *gorm.DB) error {
		for name, handler := range handlers {
			logger.Info().Str("seeder", name).Msg("Running seeder")
			if err := handler(tx); err != nil {
				return fmt.Errorf("failed to seed %s: %w", name, err)
			}
			logger.Info().Str("seeder", name).Msg("Seeder completed successfully")
		}
		return nil
	})
}

func NewSeederCommand() *cli.Command {
	injector := do.New()

	// Setup dependencies
	do.Provide(injector, injects.NewAppConfig)
	do.Provide(injector, injects.NewLogger)
	do.Provide(injector, injects.NewDatabase)

	// Seeder map
	seeder := map[string]func(tx *gorm.DB) error{
		"SeedExample": seedExample,
	}

	return &cli.Command{
		Name:  "seed",
		Usage: "Seed the database with initial data",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "all",
				Usage: "Seed all modules in one transaction",
			},
			&cli.StringFlag{
				Name:  "only",
				Usage: "Comma-separated list of specific modules to seed (e.g. --only=RootDomains,Users)",
			},
		},
		Action: func(c *cli.Context) error {
			// --all: seed everything in a transaction
			if c.Bool("all") {
				return SeedHandler(injector, seeder)
			}

			// --only: seed a list
			if only := c.String("only"); only != "" {
				names := strings.Split(only, ",")
				// Filter the map
				filteredSeeder := lo.PickByKeys(seeder, names)
				if len(filteredSeeder) == 0 {
					return fmt.Errorf("no valid seeders found for: %s", only)
				}
				return SeedHandler(injector, filteredSeeder)
			}

			return cli.ShowSubcommandHelp(c)
		},
	}
}
