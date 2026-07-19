package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

type loggerOpts struct {
	consolePretty bool
	level         slog.Level
	w             io.Writer
}

type LoggerOption interface {
	apply(lo *loggerOpts)
}

type LoggerOptions []LoggerOption

func (o LoggerOptions) apply(lo *loggerOpts) {
	for _, opt := range o {
		opt.apply(lo)
	}
}

type consolePrettyOpt bool

func (o consolePrettyOpt) apply(s *loggerOpts) { s.consolePretty = bool(o) }
func WithConsolePretty(flag bool) LoggerOption { return consolePrettyOpt(flag) }

type levelOpt slog.Level

func (o levelOpt) apply(s *loggerOpts)        { s.level = slog.Level(o) }
func WithLevel(level slog.Level) LoggerOption { return levelOpt(level) }

type writerOpt struct{ w io.Writer }

func (o writerOpt) apply(s *loggerOpts)   { s.w = o.w }
func WithWriter(w io.Writer) LoggerOption { return writerOpt{w: w} }

const (
	ansiReset              = "\u001b[0m"
	ansiFaint              = "\u001b[2m"
	ansiResetFaint         = "\u001b[22m"
	ansiBrightRed          = "\u001b[91m"
	ansiBrightGreen        = "\u001b[92m"
	ansiBrightYellow       = "\u001b[93m"
	ansiBrightRedFaint     = "\u001b[91;2m"
	ansiBrightMagenta      = "\u001b[95m"
	ansiMagenta            = "\u001b[35m"
	ansiBrightMagentaFaint = "\u001b[95;2m"
	ansiMagentaFaint       = "\u001b[35;2m"
	ansiBrightBlue         = "\u001b[94m"
	ansiBlue               = "\u001b[34m"
	ansiBrightBlueFaint    = "\u001b[94;2m"
	ansiBlueFaint          = "\u001b[34;2m"
)

// New returns the new [slog.Logger]
func New(opts ...LoggerOption) *slog.Logger {
	lopts := &loggerOpts{
		w:             os.Stdout,
		consolePretty: true,
		level:         slog.LevelDebug,
	}

	for _, o := range opts {
		o.apply(lopts)
	}

	jsonHandler := slog.NewJSONHandler(lopts.w, &slog.HandlerOptions{
		AddSource: false,
		Level:     slog.LevelDebug,
	})

	prettyHandler := tint.NewHandler(lopts.w, &tint.Options{
		Level:      lopts.level,
		AddSource:  false,
		TimeFormat: time.TimeOnly,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch level {
				case slog.LevelError:
					a.Value = slog.StringValue(fmt.Sprintf("[%s%s%s]", ansiBrightRed, "ERR", ansiReset))
				case slog.LevelWarn:
					a.Value = slog.StringValue(fmt.Sprintf("[%s%s%s]", ansiBrightYellow, "WRN", ansiReset))
				case slog.LevelInfo:
					a.Value = slog.StringValue(fmt.Sprintf("[%s%s%s]", ansiBrightGreen, "INF", ansiReset))
				case slog.LevelDebug:
					a.Value = slog.StringValue(fmt.Sprintf("[%s%s%s]", ansiMagenta, "DBG", ansiReset))
				default:
					a.Value = slog.StringValue("UNKNOWN")
				}
			}

			return a
		},
	})

	var handler slog.Handler = jsonHandler
	if lopts.consolePretty {
		handler = prettyHandler
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}
