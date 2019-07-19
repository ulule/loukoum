package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Update is the UPDATE statement.
type Update struct {
	With      With
	Table     Table
	Only      bool
	From      From
	Set       Set
	Where     Where
	Returning Returning
	Comment   Comment
}

// NewUpdate returns a new Update instance.
func NewUpdate(table Table) Update {
	return Update{
		Table: table,
		Set:   NewSet(),
	}
}

// Write exposes statement as a SQL query.
func (update Update) Write(ctx types.Context) {
	if update.IsEmpty() {
		panic("loukoum: an update statement must have a table and/or values")
	}

	if !update.With.IsEmpty() {
		update.With.Write(ctx)
		ctx.Write(" ")
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

	if !update.Comment.IsEmpty() {
		ctx.Write(token.Semicolon.String())
		ctx.Write(" ")
		update.Comment.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (update Update) IsEmpty() bool {
	return update.Table.IsEmpty() || update.Set.IsEmpty()
}

// Ensure that Update is a Statement
var _ Statement = Update{}
