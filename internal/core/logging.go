package core

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

// MustInitLogging initialize Logging and panics on error
func MustInitLogging(c AppConfig) *zerolog.Logger {
	if !c.Production {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.DurationFieldUnit = time.Microsecond
	lw := &MultiLogWriter{}
	lw.AddWriter(os.Stdout)
	if c.LogfileDir != "" {
		fw, err := os.OpenFile(c.LogfileDir+"/logfile.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(fmt.Sprintf("unable to open logfile at %s/logfile.json", c.LogfileDir))
		}
		lw.AddWriter(fw)
	}
	log.Logger = zerolog.New(lw).With().Timestamp().Logger()
	return &log.Logger
}

// LogWriter represents an io.Writer implementation that can
// be disabled in case of write error
type LogWriter struct {
	writer  io.Writer
	enabled bool
}

// MultiLogWriter is an io.Writer that can write to multiple
// LogWriter children and disable them on write error
type MultiLogWriter struct {
	writers []LogWriter
}

func (l *MultiLogWriter) AddWriter(w io.Writer) {
	l.writers = append(l.writers, LogWriter{writer: w, enabled: true})
}

func (l *MultiLogWriter) Write(p []byte) (n int, err error) {
	for _, w := range l.writers {
		if !w.enabled {
			continue
		}
		written, err := w.writer.Write(p)
		if err != nil || written != len(p) {
			w.enabled = false
			log.Error().Err(err).Msg("log writer disabled")
			continue
		}
		n = written
	}
	return n, err
}
