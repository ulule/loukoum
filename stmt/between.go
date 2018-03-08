package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Between is a BETWEEN expression.
type Between struct {
	Identifier Identifier
	Operator   ComparisonOperator
	From       Expression
	And        LogicalOperator
	To         Expression
}

// NewBetween returns a new Between instance using an inclusive operator.
func NewBetween(identifier Identifier, from, to Expression) Between {
	return Between{
		Identifier: identifier,
		Operator:   NewComparisonOperator(types.Between),
		From:       from,
		And:        NewAndOperator(),
		To:         to,
	}
}

// NewNotBetween returns a new Between instance using an exclusive operator.
func NewNotBetween(identifier Identifier, from, to Expression) Between {
	return Between{
		Identifier: identifier,
		Operator:   NewComparisonOperator(types.NotBetween),
		From:       from,
		And:        NewAndOperator(),
		To:         to,
	}
}

func (Between) expression() {}

// Write exposes statement as a SQL query.
func (between Between) Write(ctx types.Context) {
	if between.IsEmpty() {
		panic("loukoum: expression is undefined")
	}

	ctx.Write("(")
	between.Identifier.Write(ctx)
	ctx.Write(" ")
	between.Operator.Write(ctx)
	ctx.Write(" ")
	between.From.Write(ctx)
	ctx.Write(" ")
	between.And.Write(ctx)
	ctx.Write(" ")
	between.To.Write(ctx)
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (between Between) IsEmpty() bool {
	return between.Identifier.IsEmpty() || between.Operator.IsEmpty() || between.And.IsEmpty() ||
		between.From == nil || between.To == nil || between.From.IsEmpty() || between.To.IsEmpty()
}

// Ensure that Between is an Expression
var _ Expression = Between{}
