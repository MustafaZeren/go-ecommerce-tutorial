// Package logger uygulama genelinde yapılandırılmış log tutma işlemlerini sağlar.
package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New zerolog kütüphanesini kullanarak yeni bir logger örneği döner.
func New() zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	if os.Getenv("GIN_MODE") != "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
	}
	return log.Logger
}
