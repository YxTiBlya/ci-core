package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultConfig = DevelopmentConfig

var once sync.Once

func Init(cfg Config) { once.Do(func() { defaultConfig = cfg }) }

func New(name string) *Logger {
	l, err := NewCustom(name, defaultConfig)
	if err != nil {
		panic(err)
	}
	return l
}

func NewCustom(name string, cfg Config) (*Logger, error) {
	z, err := newZapWithConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Logger{z: z.Named(name)}, nil
}

type Logger struct{ z *zap.Logger }

func (l *Logger) Sync() error                       { return l.z.Sync() }
func (l *Logger) Debug(args ...interface{}) wrapper { return l.log(zap.DebugLevel, args...) }
func (l *Logger) Info(args ...interface{}) wrapper  { return l.log(zap.InfoLevel, args...) }
func (l *Logger) Warn(args ...interface{}) wrapper  { return l.log(zap.WarnLevel, args...) }
func (l *Logger) Error(args ...interface{}) wrapper { return l.log(zap.ErrorLevel, args...) }
func (l *Logger) Fatal(args ...interface{}) wrapper { return l.log(zap.FatalLevel, args...) }

func (l *Logger) log(level zapcore.Level, args ...interface{}) wrapper {
	w := wrapper{z: l.z, l: level}
	if len(args) == 0 {
		return w
	}

	w.z = w.z.With(toZapFields(w.args)...)
	return w
}
