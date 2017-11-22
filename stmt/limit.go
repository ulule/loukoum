package stmt

import (
	"strconv"

	"github.com/ulule/loukoum/types"
)

// Limit is a LIMIT clause.
type Limit struct {
	Statement
	Count int64
}

// NewLimit returns a new Limit instance.
func NewLimit(count int64) Limit {
	return Limit{
		Count: count,
	}
}

// Write expose statement as a SQL query.
func (limit Limit) Write(ctx *types.Context) {
	if limit.IsEmpty() {
		return
	}
	ctx.Write("LIMIT ")
	ctx.Write(strconv.FormatInt(limit.Count, 10))
}

// IsEmpty return true if statement is undefined.
func (limit Limit) IsEmpty() bool {
	return limit.Count == 0
}
