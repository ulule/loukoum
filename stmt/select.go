package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// SelectExpression is a SQL expression for a SELECT statement.
type SelectExpression interface {
	Statement
	selectExpression()
}

// Select is a SELECT statement.
type Select struct {
	Prefix      Prefix
	With        With
	Distinct    bool
	Expressions []SelectExpression
	From        From
	Joins       []Join
	Where       Where
	GroupBy     GroupBy
	Having      Having
	OrderBy     OrderBy
	Limit       Limit
	Offset      Offset
	Suffix      Suffix
}

// NewSelect returns a new Select instance.
func NewSelect() Select {
	return Select{}
}

// Write exposes statement as a SQL query.
func (selekt Select) Write(ctx types.Context) {
	if selekt.IsEmpty() {
		panic("loukoum: select statements must have at least one column")
	}

	selekt.writeHead(ctx)
	selekt.writeMiddle(ctx)
	selekt.writeTail(ctx)
}

func (selekt Select) writeHead(ctx types.Context) {
	if !selekt.Prefix.IsEmpty() {
		selekt.Prefix.Write(ctx)
		ctx.Write(" ")
	}

	if !selekt.With.IsEmpty() {
		selekt.With.Write(ctx)
		ctx.Write(" ")
	}

	ctx.Write(token.Select.String())

	if selekt.Distinct {
		ctx.Write(" ")
		ctx.Write(token.Distinct.String())
	}

	for i := range selekt.Expressions {
		if i == 0 {
			ctx.Write(" ")
		} else {
			ctx.Write(", ")
		}
		selekt.Expressions[i].Write(ctx)
	}

	if !selekt.From.IsEmpty() {
		ctx.Write(" ")
		selekt.From.Write(ctx)
	}
}

func (selekt Select) writeMiddle(ctx types.Context) {
	for i := range selekt.Joins {
		ctx.Write(" ")
		selekt.Joins[i].Write(ctx)
	}

	if !selekt.Where.IsEmpty() {
		ctx.Write(" ")
		selekt.Where.Write(ctx)
	}

	if !selekt.GroupBy.IsEmpty() {
		ctx.Write(" ")
		selekt.GroupBy.Write(ctx)
	}

	if !selekt.Having.IsEmpty() {
		ctx.Write(" ")
		selekt.Having.Write(ctx)
	}
}

func (selekt Select) writeTail(ctx types.Context) {
	if !selekt.OrderBy.IsEmpty() {
		ctx.Write(" ")
		selekt.OrderBy.Write(ctx)
	}

	if !selekt.Limit.IsEmpty() {
		ctx.Write(" ")
		selekt.Limit.Write(ctx)
	}

	if !selekt.Offset.IsEmpty() {
		ctx.Write(" ")
		selekt.Offset.Write(ctx)
	}

	if !selekt.Suffix.IsEmpty() {
		ctx.Write(" ")
		selekt.Suffix.Write(ctx)
	}
}

// IsEmpty returns true if statement is undefined.
func (selekt Select) IsEmpty() bool {
	return len(selekt.Expressions) == 0
}

func (Select) expression() {}

// Ensure that Select is an Expression
var _ Expression = Select{}
