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
		Pairs:  NewSetPairs(),
		List:   NewSetList(),
		IsList: false,
	}
}

// Write exposes statement as a SQL query.
func (set Set) Write(ctx *types.Context) {
	ctx.Write(token.Set.String())
	ctx.Write(" ")

	if !set.IsList && !set.Pairs.IsEmpty() {
		set.Pairs.Write(ctx)
		return
	}

	if set.IsList && !set.List.IsEmpty() {
		set.List.Write(ctx)
		return
	}

	panic("loukoum: update set must be defined")
}

// IsEmpty returns true if statement is undefined.
func (set Set) IsEmpty() bool {
	return (!set.IsList && set.Pairs.IsEmpty()) || (set.IsList && set.List.IsEmpty())
}

// Ensure that Set is a Statement.
var _ Statement = Set{}

// ----------------------------------------------------------------------------
// Pair syntax
// ----------------------------------------------------------------------------

// SetPairs is the key-value syntax (a.k.a "standard" or "default").
//
// Example:
//
//   * SET foo = 1, bar = 2, baz = 3
//
type SetPairs struct {
	Pairs map[Column]Expression
}

// NewSetPairs returns a new SetPairs instance.
func NewSetPairs() SetPairs {
	return SetPairs{
		Pairs: map[Column]Expression{},
	}
}

// Write exposes statement as a SQL query.
func (pairs SetPairs) Write(ctx *types.Context) {
	type item struct {
		k Column
		v Expression
	}

	list := make([]item, 0, len(pairs.Pairs))
	for k, v := range pairs.Pairs {
		list = append(list, item{k: k, v: v})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].k.Name < list[j].k.Name
	})

	for i := range list {
		if i != 0 {
			ctx.Write(", ")
		}

		list[i].k.Write(ctx)
		ctx.Write(" = ")
		list[i].v.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (pairs SetPairs) IsEmpty() bool {
	return len(pairs.Pairs) == 0
}

// Ensure that SetPairs is a Statement.
var _ Statement = SetPairs{}

// ----------------------------------------------------------------------------
// Column-list syntax
// ----------------------------------------------------------------------------

// SetList is the column-list syntax.
//
// Example:
//
//   * SET (foo, bar, baz) = (1, 2, 3)
//   * SET (foo, bar, baz) = (sub-select)
//
type SetList struct {
	Columns     []Column
	Expressions []Expression
}

// NewSetList returns a new SetList instance.
func NewSetList() SetList {
	return SetList{}
}

// Write exposes statement as a SQL query.
func (list SetList) Write(ctx *types.Context) {
	if list.IsEmpty() {
		panic("loukoum: set requires at least one column")
	}

	ctx.Write("(")
	for i := range list.Columns {
		if i != 0 {
			ctx.Write(", ")
		}
		list.Columns[i].Write(ctx)
	}
	ctx.Write(")")

	ctx.Write(" = ")

	ctx.Write("(")
	for i := range list.Expressions {
		if i != 0 {
			ctx.Write(", ")
		}
		list.Expressions[i].Write(ctx)
	}
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (list SetList) IsEmpty() bool {
	return len(list.Columns) == 0
}

// Ensure that SetList is a Statement.
var _ Statement = SetList{}
