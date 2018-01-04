package builder

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/ulule/loukoum/format"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
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

// rawify will replace given arguments in query to obtain a human readable statement.
// Be advised, this function is not optimized, use with caution.
func rawify(query string, args map[string]interface{}) string {
	for key, arg := range args {
		key = fmt.Sprint(":", key)
		value := format.Value(arg)
		query = strings.Replace(query, key, value, 1)
	}
	return query
}

// ToColumn takes an empty interfaces and returns a Column instance.
func ToColumn(arg interface{}) stmt.Column {
	column := stmt.Column{}

	switch value := arg.(type) {
	case string:
		column = stmt.NewColumn(value)
	case stmt.Column:
		column = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as column", arg))
	}

	if column.IsEmpty() {
		panic("loukoum: given column is undefined")
	}

	return column
}

// ToColumns takes a list of empty interfaces and returns a slice of Column instance.
func ToColumns(values []interface{}) []stmt.Column {
	// If values is a slice, we try to use recursion to obtain a slice of Column.
	if len(values) == 1 {
		switch array := values[0].(type) {
		case []stmt.Column:
			return array
		case []string:
			list := make([]interface{}, len(array))
			for i := range array {
				list[i] = array[i]
			}
			return ToColumns(list)
		}
	}

	columns := make([]stmt.Column, 0, len(values))

	for i := range values {
		columns = append(columns, ToColumn(values[i]))
	}

	return columns
}

// ToTable takes an empty interfaces and returns a Table instance.
func ToTable(arg interface{}) stmt.Table {
	table := stmt.Table{}

	switch value := arg.(type) {
	case string:
		table = stmt.NewTable(value)
	case stmt.Table:
		table = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as table", arg))
	}

	if table.IsEmpty() {
		panic("loukoum: given table is undefined")
	}

	return table
}

// ToTables takes a list of empty interfaces and returns a slice of Table instance.
func ToTables(values []interface{}) []stmt.Table {
	tables := make([]stmt.Table, 0, len(values))

	for i := range values {
		tables = append(tables, ToTable(values[i]))
	}

	return tables
}

// ToFrom takes an empty interfaces and returns a From instance.
func ToFrom(arg interface{}) stmt.From {
	from := stmt.From{}

	switch value := arg.(type) {
	case string:
		from = stmt.NewFrom(stmt.NewTable(value), false)
	case stmt.From:
		from = value
	case stmt.Table:
		from = stmt.NewFrom(value, false)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as from clause", arg))
	}

	if from.IsEmpty() {
		panic("loukoum: given from clause is undefined")
	}

	return from
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
		panic("loukoum: given suffix is undefined")
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
		panic("loukoum: given prefix is undefined")
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

// ToSetPairs takes either a types.Map or slice of types.Pair and returns
// a slice of stmt.SetPair instances.
func ToSetPairs(args []interface{}) map[stmt.Column]stmt.Expression {
	pairs := map[stmt.Column]stmt.Expression{}

	for i := range args {
		switch value := args[i].(type) {
		case map[string]interface{}:
			for k, v := range value {
				pairs[ToColumn(k)] = stmt.NewExpression(v)
			}
		case types.Map:
			for k, v := range value {
				pairs[ToColumn(k)] = stmt.NewExpression(v)
			}
		case types.Pair:
			pairs[ToColumn(value.Key)] = stmt.NewExpression(value.Value)
		}
	}

	return pairs
}

// ToSet takes either a types.Map or slice of types.Pair and returns a stmt.Set instance.
func ToSet(args []interface{}) stmt.Set {
	set := stmt.NewSet()
	pairs := ToSetPairs(args)
	for k, v := range pairs {
		set.Pairs[k] = v
	}
	return set
}
