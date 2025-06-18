package injects

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/samber/do/v2"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

func NewDatabase(i do.Injector) (*DB, error) {
	config := do.MustInvoke[ApplicationConfig](i)
	if config.MasterDB.GetDSN() == nil {
		return nil, fmt.Errorf("Master DB configuration is not set")
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  lo.Ternary(config.GoEnv == Development, logger.Info, logger.Warn),
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)

	var (
		Master *gorm.DB
		Slave  *gorm.DB
	)

	g := new(errgroup.Group)

	// Init Master DB asynchronously
	g.Go(func() error {
		db, err := gorm.Open(postgres.Open(*config.MasterDB.GetDSN()), &gorm.Config{Logger: newLogger})
		if err != nil {
			return fmt.Errorf("failed to connect to Master DB: %w", err)
		}
		Master = db
		return nil
	})

	// Init Slave DB asynchronously
	g.Go(func() error {
		var dsn string
		if config.SlaveDB.GetDSN() == nil {
			dsn = *config.MasterDB.GetDSN()
		} else {
			dsn = *config.SlaveDB.GetDSN()
		}

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})
		if err != nil {
			return fmt.Errorf("failed to connect to SlaveDB: %w", err)
		}
		Slave = db
		return nil
	})

	// Wait for both goroutines to finish
	if err := g.Wait(); err != nil {
		return nil, err
	}

	return &DB{Master: Master, Slave: Slave}, nil
}
