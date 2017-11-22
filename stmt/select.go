package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Select is a SELECT statement.
type Select struct {
	Distinct bool
	Columns  []Column
	From     From
	Joins    []Join
	Where    Where
	GroupBy  GroupBy
	Having   Having
	Limit    Limit
}

// NewSelect returns a new Select instance.
func NewSelect() Select {
	return Select{}
}

// Write expose statement as a SQL query.
func (selekt Select) Write(ctx *types.Context) {
	if selekt.IsEmpty() {
		panic("loukoum: select statements must have at least one column")
	}

	selekt.writeHead(ctx)
	selekt.writeMiddle(ctx)
	selekt.writeTail(ctx)
}

func (selekt Select) writeHead(ctx *types.Context) {
	// TODO Add prefixes

	ctx.Write("SELECT")

	if selekt.Distinct {
		ctx.Write(" DISTINCT")
	}

	for i := range selekt.Columns {
		if i == 0 {
			ctx.Write(" ")
		} else {
			ctx.Write(", ")
		}
		selekt.Columns[i].Write(ctx)
	}

	if !selekt.From.IsEmpty() {
		ctx.Write(" ")
		selekt.From.Write(ctx)
	}
}

func (selekt Select) writeMiddle(ctx *types.Context) {
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

func (selekt Select) writeTail(ctx *types.Context) {

	// TODO ORDER BY

	if !selekt.Limit.IsEmpty() {
		ctx.Write(" ")
		selekt.Limit.Write(ctx)
	}

	// TODO OFFSET

	// TODO Add suffixes
}

// IsEmpty return true if statement is undefined.
func (selekt Select) IsEmpty() bool {
	return len(selekt.Columns) == 0
}

func (Select) expression() {}
