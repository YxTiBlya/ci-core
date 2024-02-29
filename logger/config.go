package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	EncodingConsole = "console"
	EncodingJSON    = "json"
)

type Config struct {
	Level            zapcore.Level `yaml:"level" default:"-"`
	Colored          bool          `yaml:"colored"`
	TimeFormat       string        `yaml:"time_format" default:"-"`
	ConsoleSeparator string        `yaml:"console_separator"`
	Encoding         string        `yaml:"encoding"`
	Mute             bool          `yaml:"mute"`
}

var DevelopmentConfig = Config{
	Level:            zap.DebugLevel,
	Colored:          true,
	TimeFormat:       "2006-01-02T15:04:05Z07:00",
	ConsoleSeparator: "\t",
	Encoding:         EncodingConsole,
}

func newZapWithConfig(cfg Config) (*zap.Logger, error) {
	levelEncoder := zapcore.CapitalLevelEncoder
	if cfg.Colored {
		levelEncoder = zapcore.CapitalColorLevelEncoder
	}

	outputPaths := []string{"stderr"}
	errorOutputPaths := []string{"stderr"}
	if cfg.Mute {
		outputPaths = []string{}
		errorOutputPaths = []string{}
	}

	z, err := zap.Config{
		Level:         zap.NewAtomicLevelAt(cfg.Level),
		Development:   true,
		DisableCaller: true,
		Encoding:      cfg.Encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "timestamp",
			LevelKey:         "level",
			NameKey:          "name",
			CallerKey:        "caller",
			FunctionKey:      zapcore.OmitKey,
			MessageKey:       "message",
			StacktraceKey:    "stacktrace",
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      levelEncoder,
			EncodeTime:       zapcore.TimeEncoderOfLayout(cfg.TimeFormat),
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: cfg.ConsoleSeparator,
		},
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
	}.Build()
	if err != nil {
		return nil, err
	}
	return z, nil
}
