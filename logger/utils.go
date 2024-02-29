package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func toZapFields(fields []interface{}) []zap.Field {
	var zf []zap.Field
	for i := 0; i < len(fields); i += 2 {
		// If odd fields count
		if i == len(fields)-1 {
			zf = appendField(zf, "_", fields[i])
			break
		}
		zf = appendField(zf, fields[i].(string), fields[i+1])
	}
	return zf
}

func appendField(fields []zap.Field, key string, value interface{}) []zap.Field {
	switch v := value.(type) {
	case error:
		fields = append(fields, zap.Field{Key: key, Type: zapcore.StringType, String: v.Error()})
		return fields
	}
	fields = append(fields, zap.Any(key, value))
	return fields
}
