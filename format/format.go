package format

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Value formats the given value.
func Value(arg interface{}) string { // nolint: gocyclo
	switch value := arg.(type) {
	case string:
		return String(value)
	case []byte:
		return String(string(value))
	case time.Time:
		return Time(value)
	case int:
		return Int(int64(value))
	case int8:
		return Int(int64(value))
	case int16:
		return Int(int64(value))
	case int32:
		return Int(int64(value))
	case int64:
		return Int(value)
	case uint:
		return Uint(uint64(value))
	case uint8:
		return Uint(uint64(value))
	case uint16:
		return Uint(uint64(value))
	case uint32:
		return Uint(uint64(value))
	case uint64:
		return Uint(value)
	case bool:
		return Bool(value)
	case float32:
		return Float(float64(value))
	case float64:
		return Float(value)
	default:
		return fmt.Sprint(value)
	}
}

// String formats the given string.
func String(value string) string {
	var b strings.Builder
	b.WriteByte('\'')
	for _, r := range value {
		switch r {
		case '\'':
			b.WriteString(`\'`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('\'')
	return b.String()
}

// Int formats the given number.
func Int(value int64) string {
	return strconv.FormatInt(value, 10)
}

// Uint formats the given number.
func Uint(value uint64) string {
	return strconv.FormatUint(value, 10)
}

// Bool formats the given boolean.
func Bool(value bool) string {
	return strconv.FormatBool(value)
}

// Float formats the given number.
func Float(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}

// Time formats the given time.
func Time(value time.Time) string {
	return fmt.Sprint("'", value.UTC().Format("2006-01-02 15:04:05.999999"), "+00'")
}
