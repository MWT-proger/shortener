package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StringField constructs a field with the given key and value.
func StringField(key string, val string) zapcore.Field {
	return zap.String(key, val)
}

// IntField constructs a field with the given key and value.
func IntField(key string, val int) zapcore.Field {
	return zap.Int(key, val)
}

// DurationField constructs a field with the given key and value. The encoder controls how the duration is serialized.
func DurationField(key string, val time.Duration) zapcore.Field {
	return zap.Duration(key, val)
}

// ErrorField constructs a field that lazily stores err.Error() under the
// provided key. Errors which also implement fmt.Formatter (like those produced
// by github.com/pkg/errors) will also have their verbose representation stored
// under key+"Verbose". If passed a nil error, the field is a no-op.
//
// For the common case in which the key is simply "error", the Error function
// is shorter and less repetitive.
func ErrorField(err error) zapcore.Field {
	return zap.NamedError("error", err)
}
