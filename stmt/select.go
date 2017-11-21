package stmt

import (
	"github.com/ulule/loukoum/types"
)

type Select struct {
	Distinct bool
	Columns  []Column
	From     From
	Joins    []Join
	Where    Where
}

func NewSelect() Select {
	return Select{}
}

func (selekt Select) Write(ctx *types.Context) {
	if selekt.IsEmpty() {
		panic("loukoum: select statements must have at least one column")
	}

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

	for i := range selekt.Joins {
		ctx.Write(" ")
		selekt.Joins[i].Write(ctx)
	}

	if !selekt.Where.IsEmpty() {
		ctx.Write(" ")
		selekt.Where.Write(ctx)
	}

	// TODO GROUP BY

	// TODO HAVING

	// TODO ORDER BY

	// TODO LIMIT

	// TODO OFFSET

	// TODO Add suffixes
}

// IsEmpty return true if statement is undefined.
func (selekt Select) IsEmpty() bool {
	return len(selekt.Columns) == 0
}

func (selek Select) expression() {}
