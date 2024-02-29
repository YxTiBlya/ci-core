package logger

import (
	"fmt"
	"reflect"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type wrapper struct {
	z      *zap.Logger
	l      zapcore.Level
	args   []interface{}
	fields []zap.Field
}

func (w wrapper) withFields(fields ...zap.Field) wrapper {
	w.fields = append(w.fields, fields...)
	return w
}
func (w wrapper) Str(name string, value string) wrapper  { return w.withFields(zap.String(name, value)) }
func (w wrapper) Bool(name string, value bool) wrapper   { return w.withFields(zap.Bool(name, value)) }
func (w wrapper) Int(name string, value int) wrapper     { return w.withFields(zap.Int(name, value)) }
func (w wrapper) Int64(name string, value int64) wrapper { return w.withFields(zap.Int64(name, value)) }
func (w wrapper) Uint(name string, value uint) wrapper   { return w.withFields(zap.Uint(name, value)) }
func (w wrapper) Uint64(name string, value uint64) wrapper {
	return w.withFields(zap.Uint64(name, value))
}
func (w wrapper) Float64(name string, value float64) wrapper {
	return w.withFields(zap.Float64(name, value))
}
func (w wrapper) Ifc(name string, value interface{}) wrapper {
	switch reflect.TypeOf(value).Kind() {
	case reflect.Chan:
		return w.withFields(zap.String(name, "<chan>"))
	case reflect.Func:
		return w.withFields(zap.String(name, "<func>"))
	}
	return w.withFields(zap.Any(name, value))
}

func (w wrapper) Err(err error) wrapper {
	if err == nil {
		return w.withFields(zap.String("error", "<nil>"))
	}
	return w.withFields(zap.String("error", err.Error()))
}

func (w wrapper) Msg(message string, fields ...interface{}) { w.msg(message, fields...) }

func (w wrapper) Msgf(format string, args ...interface{}) { w.msg(fmt.Sprintf(format, args...)) }

func (w wrapper) msg(message string, fields ...interface{}) {
	if ce := w.z.Check(w.l, message); ce != nil {
		switch w.l {
		case zap.WarnLevel, zap.ErrorLevel, zap.FatalLevel:
			ce.Stack = ""
		}
		if len(fields) > 0 {
			ce.Write(append(w.fields, toZapFields(fields)...)...)
		} else {
			ce.Write(w.fields...)
		}
	}
}
