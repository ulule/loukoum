package stmt

import "github.com/ulule/loukoum/types"

// Into is the INTO clause.
type Into struct {
	Statement
	Table Table
}

// NewInto returns a new Into instance.
func NewInto(table Table) Into {
	return Into{
		Table: table,
	}
}

// Write implements Statement interface.
func (into Into) Write(ctx *types.Context) {
	ctx.Write("INTO ")
	into.Table.Write(ctx)
}

// IsEmpty implements Statement interface.
func (into Into) IsEmpty() bool {
	return into.Table.IsEmpty()
}
