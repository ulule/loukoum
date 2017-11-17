package builder

import (
	"fmt"
	"strings"

	"github.com/ulule/loukoum/parser"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

// Select is a builder used for "SELECT" query.
type Select struct {
	query stmt.Select
}

// NewSelect creates a new Select.
func NewSelect() Select {
	return Select{}
}

// Distinct adds a DISTINCT clause to the query.
func (b Select) Distinct() Select {
	b.query.Distinct = true
	return b
}

// Columns adds result columns to the query.
func (b Select) Columns(columns []interface{}) Select {
	if len(b.query.Columns) != 0 {
		panic("loukoum: select builder has columns already defined")
	}
	if len(columns) == 0 {
		columns = []interface{}{"*"}
	}

	b.query.Columns = ToColumns(columns)

	return b
}

// From sets the FROM clause of the query.
func (b Select) From(from interface{}) Select {
	if !b.query.From.IsEmpty() {
		panic("loukoum: select builder has from clause already defined")
	}

	switch value := from.(type) {
	case string:
		b.query.From = stmt.NewFrom(stmt.NewTable(value))
	case stmt.From:
		b.query.From = value
	case stmt.Table:
		b.query.From = stmt.NewFrom(value)
	default:
		panic(fmt.Sprintf("loukoum: cannot use %T as from clause", from))
	}

	if b.query.From.IsEmpty() {
		panic("loukoum: given from clause is undefined")
	}

	return b
}

// Join adds a JOIN clause to the query.
func (b Select) Join(args ...interface{}) Select {
	switch len(args) {
	case 1:
		return b.join1(args)
	case 2:
		return b.join2(args)
	case 3:
		return b.join3(args)
	default:
		panic("loukoum: given join clause is invalid")
	}
}

func (b Select) join1(args []interface{}) Select {
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

	b.query.Joins = append(b.query.Joins, join)

	return b
}

func (b Select) join2(args []interface{}) Select {
	join := handleSelectJoin(args)
	if join.IsEmpty() {
		panic("loukoum: given join clause is undefined")
	}

	b.query.Joins = append(b.query.Joins, join)

	return b
}

func (b Select) join3(args []interface{}) Select {
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

	b.query.Joins = append(b.query.Joins, join)

	return b
}

// Where adds WHERE clauses.
func (b Select) Where(condition stmt.Expression) Select {
	if b.query.Where.IsEmpty() {
		b.query.Where = stmt.NewWhere(condition)
		return b
	}

	return b.And(condition)
}

// And adds AND WHERE conditions.
func (b Select) And(condition stmt.Expression) Select {
	b.query.Where = b.query.Where.And(condition)
	return b
}

// Or adds OR WHERE conditions.
func (b Select) Or(condition stmt.Expression) Select {
	b.query.Where = b.query.Where.Or(condition)
	return b
}

func (b Select) String() string {
	query, args := b.Prepare()
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

func (b Select) Prepare() (string, map[string]interface{}) {
	ctx := types.NewContext()
	b.query.Write(ctx)

	query := ctx.Query()
	args := ctx.Values()

	return query, args
}

// Statement return underlying statement.
func (b Select) Statement() stmt.Statement {
	return b.query
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
