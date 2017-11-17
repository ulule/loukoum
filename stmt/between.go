package stmt

import (
	"github.com/ulule/loukoum/types"
)

type Between struct {
	Expression
	Identifier Identifier
	Operator   ComparisonOperator
	From       Expression
	And        LogicalOperator
	To         Expression
}

func NewBetween(identifier Identifier, from, to Expression) Between {
	return Between{
		Identifier: identifier,
		Operator:   NewComparisonOperator(types.Between),
		From:       from,
		And:        NewAndOperator(),
		To:         to,
	}
}

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

func (between Between) Write(ctx *types.Context) {
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

// IsEmpty return true if statement is undefined.
func (between Between) IsEmpty() bool {
	return between.Identifier.IsEmpty() || between.Operator.IsEmpty() || between.And.IsEmpty() ||
		between.From == nil || between.To == nil || between.From.IsEmpty() || between.To.IsEmpty()
}
