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

	// set log level
	level, _ := zerolog.ParseLevel(cfg.LogLevel)
	zerolog.SetGlobalLevel(level)

	// set stack trace
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// set logger output
	var loggerOutput io.Writer = zerolog.ConsoleWriter{
		Out: os.Stderr,
		// TimeFormat: time.RFC3339,
		FieldsExclude: []string{
			"http_user_agent",
			"response_size_bytes",
		},
	}

	if cfg.ENV != "development" {
		loggerOutput = os.Stderr
	}

	// set logger
	log := zerolog.New(loggerOutput).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	Log = log
}

func Fatal(err error, msg string) {
	Log.Fatal().Stack().Err(errors.Errorf("%v", err)).Msg(msg)
}
