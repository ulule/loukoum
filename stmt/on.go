package stmt

import (
	"github.com/ulule/loukoum/types"
)

// On is a ON clause.
type On struct {
	Left  Column
	Right Column
}

// NewOn returns a new On instance.
func NewOn(left, right Column) On {
	return On{
		Left:  left,
		Right: right,
	}
}

// Write expose statement as a SQL query.
func (on On) Write(ctx *types.Context) {
	ctx.Write("ON ")
	ctx.Write(on.Left.Name)
	ctx.Write(" = ")
	ctx.Write(on.Right.Name)
}

// IsEmpty return true if statement is undefined.
func (on On) IsEmpty() bool {
	return on.Left.IsEmpty() || on.Right.IsEmpty()
}
