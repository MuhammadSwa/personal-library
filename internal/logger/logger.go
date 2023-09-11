package logger

import (
	"io"
	"os"

	"github.com/muhammadswa/personal-library/config"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var Log zerolog.Logger

func init() {
	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Error initializing configuration")
	}

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var loggerOutput io.Writer = zerolog.ConsoleWriter{
		Out: os.Stderr,
		// TimeFormat: time.RFC3339,
	}
	if cfg.ENV != "development" {
		loggerOutput = os.Stderr
	}

	log := zerolog.New(loggerOutput).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	Log = log

}

// func Info(msg string, args ...any) {
// 	log.Info().Msg("hello")
// }

func Fatal(err error, msg string) {
	Log.Fatal().Stack().Err(errors.Errorf("%v", err)).Msg("Error initializing database")
}
