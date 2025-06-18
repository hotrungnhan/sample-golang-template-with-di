package injects

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
)

func NewLogger(i do.Injector) (*zerolog.Logger, error) {
	logger := zerolog.New(os.Stdout)

	return &logger, nil
}
