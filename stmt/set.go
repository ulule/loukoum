package stmt

import (
	"sort"

	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Set is a SET clause.
type Set struct {
	Values map[Column]Expression
}

// NewSet returns a new Set instance.
func NewSet() Set {
	return Set{
		Values: map[Column]Expression{},
	}
}

// Write exposes statement as a SQL query.
func (set Set) Write(ctx *types.Context) {
	ctx.Write(token.Set.String())

	type pair struct {
		k Column
		v Expression
	}

	pairs := make([]pair, 0, len(set.Values))
	for k, v := range set.Values {
		pairs = append(pairs, pair{k: k, v: v})
	}

	sort.Slice(pairs[:], func(i, j int) bool { return pairs[i].k.Name < pairs[j].k.Name })

	for i, pair := range pairs {
		if i == 0 {
			ctx.Write(" ")
		} else {
			ctx.Write(", ")
		}

		pair.k.Write(ctx)
		ctx.Write(" = ")
		pair.v.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (set Set) IsEmpty() bool {
	return len(set.Values) == 0
}

// Ensure that Set is a Statement.
var _ Statement = Set{}
