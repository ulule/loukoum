package stmt

import (
	"sort"

	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Set is a SET clause.
type Set struct {
	Pairs  SetPairs
	List   SetList
	IsList bool
}

// NewSet returns a new Set instance.
func NewSet() Set {
	return Set{
		Pairs: NewSetPairs(),
	}
}

// Write exposes statement as a SQL query.
func (set Set) Write(ctx *types.Context) {
	ctx.Write(token.Set.String())
	ctx.Write(" ")

	if !set.Pairs.IsEmpty() {
		set.Pairs.Write(ctx)
		return
	}

	if !set.List.IsEmpty() {
		set.List.Write(ctx)
		return
	}

	panic("loukoum: update set must be defined")
}

// IsEmpty returns true if statement is undefined.
func (set Set) IsEmpty() bool {
	return (!set.IsList && set.Pairs.IsEmpty()) ||
		(set.IsList && set.List.IsEmpty())
}

// Ensure that Set is a Statement.
var _ Statement = Set{}

// ----------------------------------------------------------------------------
// Pair syntax
// ----------------------------------------------------------------------------

// SetPairs is the "standard" SET key/value syntax.
type SetPairs struct {
	Pairs map[Column]Expression
}

// NewSetPairs returns a new SetPairs instance.
func NewSetPairs() SetPairs {
	return SetPairs{
		Pairs: map[Column]Expression{},
	}
}

// Merge merges new pairs into existing ones (last write wins).
func (s SetPairs) Merge(pairs map[Column]Expression) {
	for k, v := range pairs {
		s.Pairs[k] = v
	}
}

// Write exposes statement as a SQL query.
func (s SetPairs) Write(ctx *types.Context) {
	type pair struct {
		k Column
		v Expression
	}

	pairs := make([]pair, 0, len(s.Pairs))
	for k, v := range s.Pairs {
		pairs = append(pairs, pair{k: k, v: v})
	}

	sort.Slice(pairs[:], func(i, j int) bool { return pairs[i].k.Name < pairs[j].k.Name })

	for i, pair := range pairs {
		if i != 0 {
			ctx.Write(", ")
		}

		pair.k.Write(ctx)
		ctx.Write(" = ")
		pair.v.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (s SetPairs) IsEmpty() bool {
	return len(s.Pairs) == 0
}

// Ensure that SetPairs is a Statement.
var _ Statement = SetPairs{}

// ----------------------------------------------------------------------------
// Column-list syntax
// ----------------------------------------------------------------------------

// SetList is the column-list syntax.
// Example: SET (foo, bar, baz) = (1, 2, 3) / SET (foo, bar, baz) = (sub-select)
type SetList struct {
	Columns     []Column
	Expressions []Expression
}

// NewSetList returns a new SetList instance.
func NewSetList() SetList {
	return SetList{}
}

// Write exposes statement as a SQL query.
func (s SetList) Write(ctx *types.Context) {
	if s.IsEmpty() {
		panic("loukoum: set requires at least one column")
	}

	ctx.Write("(")
	for i := range s.Columns {
		if i != 0 {
			ctx.Write(", ")
		}
		s.Columns[i].Write(ctx)
	}
	ctx.Write(")")

	ctx.Write(" = ")

	ctx.Write("(")
	for i := range s.Expressions {
		if i != 0 {
			ctx.Write(", ")
		}
		s.Expressions[i].Write(ctx)
	}
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (s SetList) IsEmpty() bool {
	return len(s.Columns) == 0
}

// Ensure that SetList is a Statement.
var _ Statement = SetList{}
