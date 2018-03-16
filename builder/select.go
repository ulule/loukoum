package builder

import (
	"fmt"

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
func (b Select) Columns(args ...interface{}) Select {
	if len(b.query.Columns) != 0 {
		panic("loukoum: select builder has columns already defined")
	}
	if len(args) == 0 {
		args = []interface{}{"*"}
	}

	b.query.Columns = ToColumns(args)

	return b
}

// From sets the FROM clause of the query.
func (b Select) From(arg interface{}) Select {
	if !b.query.From.IsEmpty() {
		panic("loukoum: select builder has from clause already defined")
	}

	b.query.From = ToFrom(arg)

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

// GroupBy adds GROUP BY clauses.
func (b Select) GroupBy(args ...interface{}) Select {
	if !b.query.GroupBy.IsEmpty() {
		panic("loukoum: select builder has group by clause already defined")
	}

	columns := ToColumns(args)
	group := stmt.NewGroupBy(columns)
	if group.IsEmpty() {
		panic("loukoum: given join clause is undefined")
	}

	b.query.GroupBy = group

	return b
}

// Having adds HAVING clauses.
func (b Select) Having(condition stmt.Expression) Select {
	if !b.query.Having.IsEmpty() {
		panic("loukoum: select builder has having clause already defined")
	}

	b.query.Having = stmt.NewHaving(condition)

	return b
}

// OrderBy adds ORDER BY clauses.
func (b Select) OrderBy(orders ...stmt.Order) Select {
	b.query.OrderBy.Orders = append(b.query.OrderBy.Orders, orders...)
	return b
}

// Limit adds LIMIT clause.
func (b Select) Limit(value interface{}) Select {
	if !b.query.Limit.IsEmpty() {
		panic("loukoum: select builder has limit clause already defined")
	}

	limit, ok := ToInt64(value)
	if !ok || limit <= 0 {
		panic("loukoum: limit must be a positive integer")
	}

	b.query.Limit = stmt.NewLimit(limit)

	return b
}

// Offset adds OFFSET clause.
func (b Select) Offset(value interface{}) Select {
	if !b.query.Offset.IsEmpty() {
		panic("loukoum: select builder has offset clause already defined")
	}

	offset, ok := ToInt64(value)
	if !ok || offset < 0 {
		panic("loukoum: offset must be a non-negative integer")
	}

	b.query.Offset = stmt.NewOffset(offset)

	return b
}

// Suffix adds given clauses as suffixes.
func (b Select) Suffix(suffix interface{}) Select {
	if !b.query.Offset.IsEmpty() {
		panic("loukoum: select builder has suffixes clauses already defined")
	}

	b.query.Suffix = ToSuffix(suffix)

	return b
}

// Prefix adds given clauses as prefixes.
func (b Select) Prefix(prefix interface{}) Select {
	if !b.query.Offset.IsEmpty() {
		panic("loukoum: select builder has prefixes clauses already defined")
	}

	b.query.Prefix = ToPrefix(prefix)

	return b
}

// String returns the underlying query as a raw statement.
// This function should be used for debugging since it doesn't escape anything and is completely
// vulnerable to SQL injection.
// You should use either NamedQuery() or Query()...
func (b Select) String() string {
	var ctx types.RawContext
	b.query.Write(&ctx)
	return ctx.Query()
}

// NamedQuery returns the underlying query as a named statement.
func (b Select) NamedQuery() (string, map[string]interface{}) {
	var ctx types.NamedContext
	b.query.Write(&ctx)
	return ctx.Query(), ctx.Values()
}

// Query returns the underlying query as a regular statement.
func (b Select) Query() (string, []interface{}) {
	var ctx types.StdContext
	b.query.Write(&ctx)
	return ctx.Query(), ctx.Values()
}

// Statement returns underlying statement.
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

// Ensure that Select is a Builder
var _ Builder = Select{}
