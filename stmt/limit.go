package stmt

import (
	"strconv"

	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// Limit is a LIMIT clause.
type Limit struct {
	Count int64
}

// NewLimit returns a new Limit instance.
func NewLimit(count int64) Limit {
	return Limit{
		Count: count,
	}
}

// Write exposes statement as a SQL query.
func (limit Limit) Write(ctx types.Context) {
	if limit.IsEmpty() {
		return
	}
	ctx.Write(token.Limit.String())
	ctx.Write(" ")
	ctx.Write(strconv.FormatInt(limit.Count, 10))
}

// IsEmpty returns true if statement is undefined.
func (limit Limit) IsEmpty() bool {
	return limit.Count == 0
}

// Ensure that Limit is a Statement
var _ Statement = Limit{}
