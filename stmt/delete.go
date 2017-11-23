package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Delete is a DELETE statement.
type Delete struct {
	From      From
	Using     Using
	Where     Where
	Returning Returning
}

// NewDelete returns a new Delete instance.
func NewDelete() Delete {
	return Delete{}
}

// Write exposes statement as a SQL query.
func (delete Delete) Write(ctx *types.Context) {
	if delete.IsEmpty() {
		panic("loukoum: a delete statement must have a table")
	}

	ctx.Write(token.Delete.String())
	ctx.Write(" ")
	delete.From.Write(ctx)

	if !delete.Using.IsEmpty() {
		ctx.Write(" ")
		delete.Using.Write(ctx)
	}

	if !delete.Where.IsEmpty() {
		ctx.Write(" ")
		delete.Where.Write(ctx)
	}

	if !delete.Returning.IsEmpty() {
		ctx.Write(" ")
		delete.Returning.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (delete Delete) IsEmpty() bool {
	return delete.From.IsEmpty()
}

// Ensure that Delete is a Statement
var _ Statement = Delete{}
