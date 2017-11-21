package stmt

import (
	"github.com/ulule/loukoum/types"
)

type In struct {
	Expression
	Identifier Identifier
	Operator   ComparisonOperator
	Value      Expression
}

func NewIn(identifier Identifier, value Expression) In {
	return In{
		Identifier: identifier,
		Operator:   NewComparisonOperator(types.In),
		Value:      value,
	}
}

func NewNotIn(identifier Identifier, value Expression) In {
	return In{
		Identifier: identifier,
		Operator:   NewComparisonOperator(types.NotIn),
		Value:      value,
	}
}

func (In) expression() {}

func (in In) Write(ctx *types.Context) {
	if in.IsEmpty() {
		panic("loukoum: expression is undefined")
	}

	ctx.Write("(")
	in.Identifier.Write(ctx)
	ctx.Write(" ")
	in.Operator.Write(ctx)
	ctx.Write(" (")
	in.Value.Write(ctx)
	ctx.Write("))")
}

// IsEmpty return true if statement is undefined.
func (in In) IsEmpty() bool {
	return in.Identifier.IsEmpty() || in.Operator.IsEmpty() || in.Value == nil || in.Value.IsEmpty()
}
