package stmt

import (
	"sort"

	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Set is a SET clause.
type Set struct {
	Pairs PairContainer
}

// NewSet returns a new Set instance.
func NewSet() Set {
	return Set{
		Pairs: NewPairContainer(),
	}
}

// Write exposes statement as a SQL query.
func (set Set) Write(ctx *types.Context) {
	ctx.Write(token.Set.String())
	ctx.Write(" ")
	set.Pairs.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (set Set) IsEmpty() bool {
	return set.Pairs.IsEmpty()
}

// Ensure that Set is a Statement.
var _ Statement = Set{}

// PairMode define the mode of PairContainer.
type PairMode uint8

const (
	// PairUnknownMode define an unknown mode.
	PairUnknownMode = PairMode(iota)
	// PairAssociativeMode define a key-value mode for PairContainer.
	PairAssociativeMode
	// PairArrayMode define a column-list mode for PairContainer.
	PairArrayMode
)

// PairContainer is a composite collection that store a list of values for SET clause.
type PairContainer struct {
	Mode        PairMode
	Map         map[Column]Expression
	Columns     []Column
	Expressions []Expression
}

// NewPairContainer creates a new PairContainer.
func NewPairContainer() PairContainer {
	return PairContainer{
		Mode:        PairUnknownMode,
		Map:         map[Column]Expression{},
		Columns:     []Column{},
		Expressions: []Expression{},
	}
}

// Add appends given column and expression.
// It will configure Set's syntax to key-value (a.k.a "standard", "default" or "associative").
//
// Example:
//
//   * SET foo = 1, bar = 2, baz = 3
//
func (pairs *PairContainer) Add(column Column, expression Expression) {
	if pairs.Mode == PairUnknownMode {
		pairs.Mode = PairAssociativeMode
	}
	if pairs.Mode != PairAssociativeMode {
		panic("loukoum: you can only use pairs in key-value or column-list syntax")
	}

	_, ok := pairs.Map[column]
	if !ok {
		pairs.Columns = append(pairs.Columns, column)
	}

	pairs.Map[column] = expression
}

// Set appends given column. It will configure Set's syntax to column-list.
// You may use Use(...) function to provide required expressions.
//
// Example:
//
//   * SET (foo, bar, baz) = (1, 2, 3)
//   * SET (foo, bar, baz) = (sub-select)
//
func (pairs *PairContainer) Set(column Column) {
	if pairs.Mode == PairUnknownMode {
		pairs.Mode = PairArrayMode
	}
	if pairs.Mode != PairArrayMode {
		panic("loukoum: you can only use pairs in key-value or column-list syntax")
	}

	pairs.Columns = append(pairs.Columns, column)
}

// Use appends given expression if Set's syntax is defined column-list.
// You have to use Set(...) function to provide required columns.
//
// Example:
//
//   * SET (foo, bar, baz) = (1, 2, 3)
//   * SET (foo, bar, baz) = (sub-select)
//
func (pairs *PairContainer) Use(expression Expression) {
	if pairs.Mode != PairArrayMode {
		panic("loukoum: you have to define pairs columns first")
	}

	pairs.Expressions = append(pairs.Expressions, expression)
}

// Values returns columns and expressions of current instance.
func (pairs PairContainer) Values() ([]Column, []Expression) {
	if pairs.Mode != PairAssociativeMode {
		return pairs.Columns, pairs.Expressions
	}

	sort.Slice(pairs.Columns, func(i, j int) bool {
		return pairs.Columns[i].Name < pairs.Columns[j].Name ||
			pairs.Columns[i].Alias < pairs.Columns[j].Alias
	})

	columns := make([]Column, 0, len(pairs.Columns))
	expressions := make([]Expression, 0, len(pairs.Columns))

	for i := range pairs.Columns {
		column := pairs.Columns[i]
		expression, ok := pairs.Map[column]
		if !ok {
			panic("loukoum: invalid state for stmt.PairContainer")
		}
		columns = append(columns, column)
		expressions = append(expressions, expression)
	}

	return columns, expressions
}

// Write exposes statement as a SQL query.
func (pairs PairContainer) Write(ctx *types.Context) {
	if pairs.IsEmpty() {
		panic("loukoum: values for SET clause are required")
	}
	if pairs.Mode == PairAssociativeMode {
		pairs.WriteAssociative(ctx)
	}
	if pairs.Mode == PairArrayMode {
		pairs.WriteArray(ctx)
	}
}

// WriteAssociative exposes statement as a SQL query using a key-value syntax.
func (pairs PairContainer) WriteAssociative(ctx *types.Context) {
	columns, expressions := pairs.Values()

	for i := range columns {
		if i != 0 {
			ctx.Write(", ")
		}

		columns[i].Write(ctx)
		ctx.Write(" = ")
		expressions[i].Write(ctx)
	}
}

// WriteArray exposes statement as a SQL query using a column-list syntax.
func (pairs PairContainer) WriteArray(ctx *types.Context) {
	ctx.Write("(")
	for i := range pairs.Columns {
		if i != 0 {
			ctx.Write(", ")
		}
		pairs.Columns[i].Write(ctx)
	}
	ctx.Write(")")

	ctx.Write(" = ")

	ctx.Write("(")
	for i := range pairs.Expressions {
		if i != 0 {
			ctx.Write(", ")
		}
		pairs.Expressions[i].Write(ctx)
	}
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (pairs PairContainer) IsEmpty() bool {
	return pairs.Mode == PairUnknownMode || (len(pairs.Map) == 0 && len(pairs.Expressions) == 0)
}

// Ensure that PairContainer is a Statement.
var _ Statement = PairContainer{}
