package logger

import (
	"os"

	"github.com/rs/zerolog"
)

func NewLogger() zerolog.Logger {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	logger.Info().Msg("Zerolog logger initialised")
	return logger
}
