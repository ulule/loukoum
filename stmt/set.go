package stmt

import (
	"sort"

	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// SetPair is a set pair.
type SetPair struct {
	Column     Column
	Expression Expression
}

// NewSetPair returns a new SetPair instance.
func NewSetPair(column Column, expression Expression) SetPair {
	return SetPair{
		Column:     column,
		Expression: expression,
	}
}

// Set is a SET clause.
type Set struct {
	Pairs []SetPair
}

// NewSet returns a new Set instance.
func NewSet() Set {
	return Set{}
}

// Write exposes statement as a SQL query.
func (set Set) Write(ctx *types.Context) {
	ctx.Write(token.Set.String())

	sort.Slice(set.Pairs[:], func(i, j int) bool { return set.Pairs[i].Column.Name < set.Pairs[j].Column.Name })

	for i, pair := range set.Pairs {
		if i == 0 {
			ctx.Write(" ")
		} else {
			ctx.Write(", ")
		}

		pair.Column.Write(ctx)
		ctx.Write(" = ")
		pair.Expression.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (set Set) IsEmpty() bool {
	return len(set.Pairs) == 0
}

// Ensure that Set is a Statement.
var _ Statement = Set{}
