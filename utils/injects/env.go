package injects

import (
	"fmt"
	"github.com/hotrungnhan/surl/utils/helpers"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/samber/do/v2"
	"github.com/spf13/viper"
)

type DBConfig struct {
	Host     string `mapstructure:"HOST" default:"localhost"`
	Port     int    `mapstructure:"PORT" default:"5432"`
	Username string `mapstructure:"USERNAME" default:"postgres"`
	Password string `mapstructure:"PASSWORD" default:"postgres"`
	Name     string `mapstructure:"NAME" default:"example_db"`  //
	SslMode  string `mapstructure:"SSL_MODE" default:"disable"` // e.g., "disable", "require", "verify-ca", "verify-full"
}

func (dc DBConfig) GetDSN() *string {
	if dc.Host == "" && dc.Port == 0 && dc.Username == "" && dc.Password == "" && dc.Name == "" {
		return nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		dc.Host,
		dc.Username,
		dc.Password,
		dc.Name,
		fmt.Sprintf("%d", dc.Port),
		dc.SslMode,
		"UTC",
	)

	return &dsn
}

type CacheDBConfig struct {
}

type FeatureFlags struct {
	ENABLE_FEATURE_X bool `mapstructure:"ENABLE_FEATURE_X"` // Example feature flag
}

type GoEnv string

var (
	Development GoEnv = "development"
	Production  GoEnv = "production"
	Testing     GoEnv = "testing"
	Staging     GoEnv = "staging"
)

func (g *GoEnv) Short() string {
	switch *g {
	case Development:
		return "dev"
	case Production:
		return "prod"
	case Testing:
		return "test"
	case Staging:
		return "stg"
	default:
		return "unknown	"
	}
}

type ApplicationConfig struct {
	MasterDB DBConfig      `mapstructure:"MASTER_DB"`
	SlaveDB  DBConfig      `mapstructure:"SLAVE_DB"`
	CacheDB  CacheDBConfig `mapstructure:"CACHE_DB"`

	TimeZone string `mapstructure:"TIME_ZONE"` // e.g., "UTC", "America/New_York"

	LogLevel string `mapstructure:"LOG_LEVEL"` // e.g., "debug", "info", "warn", "error"

	GoEnv        GoEnv        `mapstructure:"GO_ENV" default:"development" validate:"oneof=development production testing"`
	FeatureFlags FeatureFlags `mapstructure:"FEATURE_FLAGS" default:"-"` // e.g., "new_ui": "true", "beta_feature": "false"
}

func NewAppConfig(i do.Injector) (ApplicationConfig, error) {
	// Load .env file into os environment
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or error loading it")
	}

	v := viper.NewWithOptions(viper.ExperimentalBindStruct())

	// Tell Viper to read env variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv() // Automatically read environment variables
	var cfg ApplicationConfig

	err = helpers.SetDefaults(&cfg)

	if err != nil {
		return cfg, fmt.Errorf("failed to set defaults: %w", err)
	}

	err = helpers.Validate(&cfg)

	if err != nil {
		return cfg, fmt.Errorf("validation failed: %w", err)
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
