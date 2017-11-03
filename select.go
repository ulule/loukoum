package loukoum

import (
	"bytes"
	"fmt"

	"github.com/ulule/loukoum/parser"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

// SelectBuilder is a builder used for "SELECT" query.
type SelectBuilder struct {
	distinct bool
	columns  []stmt.Column
	from     stmt.From
	joins    []stmt.Join
}

// NewSelectBuilder creates a new SelectBuilder.
func NewSelectBuilder() SelectBuilder {
	return SelectBuilder{}
}

// Distinct adds a DISTINCT clause to the query.
func (builder SelectBuilder) Distinct() SelectBuilder {
	builder.distinct = true
	return builder
}

// Columns adds result columns to the query.
func (builder SelectBuilder) Columns(columns []interface{}) SelectBuilder {
	if len(builder.columns) != 0 {
		panic("loukoum: select builder has columns already defined")
	}
	if len(columns) == 0 {
		columns = []interface{}{"*"}
	}

	builder.columns = []stmt.Column{}
	for i := range columns {
		column := stmt.Column{}
		switch value := columns[i].(type) {
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
		builder.columns = append(builder.columns, column)
	}

	return builder
}

// From sets the FROM clause of the query.
func (builder SelectBuilder) From(from interface{}) SelectBuilder {
	if !builder.from.IsEmpty() {
		panic("loukoum: select builder has from clause already defined")
	}

	switch value := from.(type) {
	case string:
		builder.from = stmt.NewFrom(stmt.NewTable(value))
	case stmt.From:
		builder.from = value
	case stmt.Table:
		builder.from = stmt.NewFrom(value)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as from clause", from))
	}

	if builder.from.IsEmpty() {
		panic("loukoum: given from clause is undefined")
	}

	return builder
}

// Join adds a JOIN clause to the query.
func (builder SelectBuilder) Join(args ...interface{}) SelectBuilder {
	switch len(args) {
	case 1:
		return builder.join1(args)
	case 2:
		return builder.join2(args)
	case 3:
		return builder.join3(args)
	default:
		panic("loukoum: given join clause is invalid")
	}
}

func (builder SelectBuilder) join1(args []interface{}) SelectBuilder {
	join := stmt.Join{}

	switch value := args[0].(type) {
	case string:
		join = parser.MustParseJoin(value)
	case stmt.Join:
		join = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as join clause", args[0]))
	}

	if join.IsEmpty() {
		panic("loukoum: given join clause is undefined")
	}

	builder.joins = append(builder.joins, join)
	return builder
}

func (builder SelectBuilder) join2(args []interface{}) SelectBuilder {
	join := handleSelectJoin(args)
	if join.IsEmpty() {
		panic("loukoum: given join clause is undefined")
	}

	builder.joins = append(builder.joins, join)
	return builder
}

func (builder SelectBuilder) join3(args []interface{}) SelectBuilder {
	join := handleSelectJoin(args)

	switch value := args[2].(type) {
	case types.JoinType:
		join.Type = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as join clause", args[1]))
	}

	if join.IsEmpty() {
		panic("loukoum: given join clause is undefined")
	}

	builder.joins = append(builder.joins, join)
	return builder
}

func handleSelectJoin(args []interface{}) stmt.Join {
	join := stmt.Join{}
	table := stmt.Table{}

	switch value := args[0].(type) {
	case string:
		table = stmt.NewTable(value)
	case stmt.Table:
		table = value
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as table argument for join clause", args[0]))
	}

	switch value := args[1].(type) {
	case string:
		join = parser.MustParseJoin(value)
	case stmt.On:
		join = stmt.NewInnerJoin(table, value)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as condition for join clause", args[1]))
	}

	join.Table = table
	return join
}

func (builder SelectBuilder) String() string {
	if len(builder.columns) == 0 {
		panic("loukoum: select statements must have at least one column")
	}

	buffer := &bytes.Buffer{}

	// TODO Add prefixes

	buffer.WriteString("SELECT ")

	if builder.distinct {
		buffer.WriteString("DISTINCT ")
	}

	for i := range builder.columns {
		if i != 0 {
			buffer.WriteString(", ")
		}
		builder.columns[i].Write(buffer)
	}

	if !builder.from.IsEmpty() {
		buffer.WriteString(" ")
		builder.from.Write(buffer)
	}

	for i := range builder.joins {
		buffer.WriteString(" ")
		builder.joins[i].Write(buffer)
	}

	// TODO WHERE

	// TODO GROUP BY

	// TODO HAVING

	// TODO ORDER BY

	// TODO LIMIT

	// TODO OFFSET

	// TODO Add suffixes

	return buffer.String()
}
