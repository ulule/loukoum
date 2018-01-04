package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Insert is a INSERT statement.
type Insert struct {
	Into       Into
	Columns    []Column
	Values     Values
	OnConflict OnConflict
	Returning  Returning
}

// NewInsert returns a new Insert instance.
func NewInsert() Insert {
	return Insert{}
}

// Write exposes statement as a SQL query.
func (insert Insert) Write(ctx *types.Context) {
	if insert.IsEmpty() {
		panic("loukoum: an insert statement must have at least one column")
	}

	ctx.Write(token.Insert.String())
	ctx.Write(" ")
	insert.Into.Write(ctx)

	if len(insert.Columns) > 0 {
		nbColumns := len(insert.Columns)
		for i := range insert.Columns {
			if i == 0 {
				ctx.Write(" (")
			} else {
				ctx.Write(", ")
			}
			insert.Columns[i].Write(ctx)
			if i == nbColumns-1 {
				ctx.Write(")")
			}
		}
	}

	if !insert.Values.IsEmpty() {
		ctx.Write(" ")
		insert.Values.Write(ctx)
	}

	if !insert.OnConflict.IsEmpty() {
		ctx.Write(" ")
		insert.OnConflict.Write(ctx)
	}

	if !insert.Returning.IsEmpty() {
		ctx.Write(" ")
		insert.Returning.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (insert Insert) IsEmpty() bool {
	return insert.Into.IsEmpty()
}

// Ensure that Insert is a Statement
var _ Statement = Insert{}
