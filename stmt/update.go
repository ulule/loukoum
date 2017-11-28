package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Update is the UPDATE statement.
type Update struct {
	Table     Table
	Only      bool
	From      From
	Set       Set
	Where     Where
	Returning Returning
}

// NewUpdate returns a new Update instance.
func NewUpdate(table Table) Update {
	return Update{
		Table: table,
	}
}

// Write exposes statement as a SQL query.
func (update Update) Write(ctx *types.Context) {
	if update.IsEmpty() {
		panic("loukoum: an update statement must have a table")
	}

	ctx.Write(token.Update.String())

	if update.Only {
		ctx.Write(" ")
		ctx.Write(token.Only.String())
	}

	ctx.Write(" ")
	update.Table.Write(ctx)

	ctx.Write(" ")
	update.Set.Write(ctx)

	if !update.From.IsEmpty() {
		ctx.Write(" ")
		update.From.Write(ctx)
	}

	if !update.Where.IsEmpty() {
		ctx.Write(" ")
		update.Where.Write(ctx)
	}

	if !update.Returning.IsEmpty() {
		ctx.Write(" ")
		update.Returning.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (update Update) IsEmpty() bool {
	return update.Table.IsEmpty() || update.Set.IsEmpty()
}

// Ensure that Update is a Statement
var _ Statement = Update{}
