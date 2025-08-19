package logging
// Logger initialization placeholder
import (
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/mdeadwioer/EyeSky/internal/platform/config"
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