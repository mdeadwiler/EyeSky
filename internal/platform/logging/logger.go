package logging
// Logger initialization placeholder
import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"

	"github.com/mdeadwiler/EyeSky/internal/platform/config"
)

type Logger struct {
	logger zerolog.Logger
}


func New(cfg config.LoggingConfig) *Logger {
 // output writer
 var output io.Writer = os.Stdout

 if cfg.Format == "console" {
	output = zerolog.ConsoleWriter{Out: os.Stdout}
 }

 // Parsing and set the level of logging
 level, err := zerolog.ParseLevel(strings.ToLower(cfg.Level))
 if err != nil {
	level = zerolog.InfoLevel // This is fallback info
 }

 // Logger
 zl := zerolog.New(output).
 Level(level).
 With().
 Timestamp().
 Logger()
 return &Logger{logger:zl}
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
 }

func (l *Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}

