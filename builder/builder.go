package builder

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ulule/loukoum/stmt"
)

// Builder defines a generic methods available for Select, Insert, Update and Delete builders.
type Builder interface {
	// String returns the underlying query as a raw statement.
	String() string
	// Prepare returns the underlying query as a named statement.
	Prepare() (string, map[string]interface{})
	// Statement returns underlying statement.
	Statement() stmt.Statement
}

func rawify(query string, args map[string]interface{}) string {
	for key, arg := range args {
		switch value := arg.(type) {
		case string:
			query = strings.Replace(query, key, fmt.Sprint("'", value, "'"), 1)
		default:
			query = strings.Replace(query, key, fmt.Sprint(value), 1)
		}
	}
	return query
}

// ToColumns takes a list of empty interfaces and returns a slice of Column instance.
func ToColumns(values []interface{}) []stmt.Column {
	columns := make([]stmt.Column, 0, len(values))

	for i := range values {
		column := stmt.Column{}

		switch value := values[i].(type) {
		case string:
			column = stmt.NewColumn(value)
		case stmt.Column:
			column = value
		default:
			panic(fmt.Sprintf("loukoum: cannot use %T as column", column))
		}
		if column.IsEmpty() {
			panic("loukoum: a column was undefined")
		}

		columns = append(columns, column)
	}

	return columns
}

// ToSuffix takes an empty interfaces and returns a Suffix instance.
func ToSuffix(arg interface{}) stmt.Suffix {
	suffix := stmt.Suffix{}

	switch value := arg.(type) {
	case string:
		suffix = stmt.NewSuffix(value)
	case stmt.Suffix:
		suffix = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as suffix", value))
	}

	if suffix.IsEmpty() {
		panic("loukoum: a suffix was undefined")
	}

	return suffix
}

// ToPrefix takes an empty interfaces and returns a Prefix instance.
func ToPrefix(arg interface{}) stmt.Prefix {
	prefix := stmt.Prefix{}

	switch value := arg.(type) {
	case string:
		prefix = stmt.NewPrefix(value)
	case stmt.Prefix:
		prefix = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as prefix", value))
	}

	if prefix.IsEmpty() {
		panic("loukoum: a prefix was undefined")
	}

	return prefix
}

// ToInt64 takes an empty interfaces and returns a int64.
func ToInt64(value interface{}) (int64, bool) { // nolint: gocyclo
	switch cast := value.(type) {
	case int64:
		return cast, true
	case int:
		return int64(cast), true
	case int8:
		return int64(cast), true
	case int16:
		return int64(cast), true
	case int32:
		return int64(cast), true
	case uint8:
		return int64(cast), true
	case uint16:
		return int64(cast), true
	case uint32:
		return int64(cast), true
	case uint64:
		if cast <= math.MaxInt64 {
			return int64(cast), true
		}
	case string:
		n, err := strconv.ParseInt(cast, 10, 64)
		if err == nil {
			return n, true
		}
	}
	return 0, false
}
