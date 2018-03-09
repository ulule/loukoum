package stmt

import (
	"strconv"

	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Offset is a OFFSET clause.
type Offset struct {
	Start int64
}

// NewOffset returns a new Offset instance.
func NewOffset(start int64) Offset {
	return Offset{
		Start: start,
	}
}

// Write exposes statement as a SQL query.
func (offset Offset) Write(ctx types.Context) {
	if offset.IsEmpty() {
		return
	}
	ctx.Write(token.Offset.String())
	ctx.Write(" ")
	ctx.Write(strconv.FormatInt(offset.Start, 10))
}

// IsEmpty returns true if statement is undefined.
func (offset Offset) IsEmpty() bool {
	return offset.Start == 0
}

// Ensure that Offset is a Statement
var _ Statement = Offset{}
