package loukoum

import (
	"fmt"
	"strings"

	"github.com/ulule/loukoum/parser"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

// SelectBuilder is a builder used for "SELECT" query.
type SelectBuilder struct {
	query stmt.Select
}

// NewSelectBuilder creates a new SelectBuilder.
func NewSelectBuilder() SelectBuilder {
	return SelectBuilder{}
}

// Distinct adds a DISTINCT clause to the query.
func (builder SelectBuilder) Distinct() SelectBuilder {
	builder.query.Distinct = true
	return builder
}

// Columns adds result columns to the query.
func (builder SelectBuilder) Columns(columns []interface{}) SelectBuilder {
	if len(builder.query.Columns) != 0 {
		panic("loukoum: select builder has columns already defined")
	}
	if len(columns) == 0 {
		columns = []interface{}{"*"}
	}

	builder.query.Columns = make([]stmt.Column, 0, len(columns))
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
		builder.query.Columns = append(builder.query.Columns, column)
	}

	return builder
}

// From sets the FROM clause of the query.
func (builder SelectBuilder) From(from interface{}) SelectBuilder {
	if !builder.query.From.IsEmpty() {
		panic("loukoum: select builder has from clause already defined")
	}

	switch value := from.(type) {
	case string:
		builder.query.From = stmt.NewFrom(stmt.NewTable(value))
	case stmt.From:
		builder.query.From = value
	case stmt.Table:
		builder.query.From = stmt.NewFrom(value)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as from clause", from))
	}

	if builder.query.From.IsEmpty() {
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

	builder.query.Joins = append(builder.query.Joins, join)
	return builder
}

func (builder SelectBuilder) join2(args []interface{}) SelectBuilder {
	join := handleSelectJoin(args)
	if join.IsEmpty() {
		panic("loukoum: given join clause is undefined")
	}

	builder.query.Joins = append(builder.query.Joins, join)
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

	builder.query.Joins = append(builder.query.Joins, join)
	return builder
}

// Where adds WHERE clauses.
func (builder SelectBuilder) Where(condition stmt.Expression) SelectBuilder {
	if builder.query.Where.IsEmpty() {
		builder.query.Where = stmt.NewWhere(condition)
		return builder
	}

	return builder.And(condition)
}

// And adds AND WHERE conditions.
func (builder SelectBuilder) And(condition stmt.Expression) SelectBuilder {
	builder.query.Where = builder.query.Where.And(condition)
	return builder
}

// Or adds OR WHERE conditions.
func (builder SelectBuilder) Or(condition stmt.Expression) SelectBuilder {
	builder.query.Where = builder.query.Where.Or(condition)
	return builder
}

func (builder SelectBuilder) String() string {
	query, args := builder.Prepare()
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

func (builder SelectBuilder) Prepare() (string, map[string]interface{}) {
	ctx := types.NewContext()
	builder.query.Write(ctx)

	query := ctx.Query()
	args := ctx.Values()

	return query, args
}

// Statement return underlying statement.
func (builder SelectBuilder) Statement() stmt.Statement {
	return builder.query
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
