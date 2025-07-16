// Package logging is to initialize the logger package
package logging

import (
	"os"

	"github.com/rs/zerolog"
)

var Logger zerolog.Logger

func Init(debug bool) {
	level := zerolog.InfoLevel
	if debug {
		level = zerolog.DebugLevel
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Logger = zerolog.New(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "2006-01-02 15:04:05",
	}).Level(level).With().Timestamp().Logger()
}

func NewContextLogger(context string) zerolog.Logger {
	return Logger.With().Str("context", context).Logger()
}
