package format

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

// Value formats the given value.
func Value(arg interface{}) string { // nolint: gocyclo
	if arg == nil {
		return "NULL"
	}

	switch value := arg.(type) {
	case string:
		return String(value)
	case []byte:
		return Bytes(value)
	case time.Time:
		return Time(value)
	case driver.Valuer:
		v, err := value.Value()
		if err != nil {
			panic("loukoum: was not able to retrieve valuer value")
		}
		return Value(v)
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
	buffer := &bytes.Buffer{}
	writeRune(buffer, '\'')
	for _, char := range value {
		switch char {
		case '\'':
			writeString(buffer, `\'`)
		case '\\':
			writeString(buffer, `\\`)
		case '\n':
			writeString(buffer, `\n`)
		case '\r':
			writeString(buffer, `\r`)
		case '\t':
			writeString(buffer, `\t`)
		default:
			writeRune(buffer, char)
		}
	}
	writeRune(buffer, '\'')
	return buffer.String()
}

// Bytes formats the give bytes.
func Bytes(value []byte) string {
	encoded := hex.EncodeToString(value)
	return fmt.Sprintf("decode('%s', 'hex')", encoded)
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

func writeRune(buffer *bytes.Buffer, chunk rune) {
	_, err := buffer.WriteRune(chunk)
	if err != nil {
		panic("loukoum: cannot write on bytes buffer")
	}
}

func writeString(buffer *bytes.Buffer, chunk string) {
	_, err := buffer.WriteString(chunk)
	if err != nil {
		panic("loukoum: cannot write on bytes buffer")
	}
}
