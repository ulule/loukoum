package stmt

import (
	"bytes"

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

func (between Between) Write(buffer *bytes.Buffer) {
	if between.IsEmpty() {
		panic("loukoum: expression is undefined")
	}

	buffer.WriteString("(")
	between.Identifier.Write(buffer)
	buffer.WriteString(" ")
	between.Operator.Write(buffer)
	buffer.WriteString(" ")
	between.From.Write(buffer)
	buffer.WriteString(" ")
	between.And.Write(buffer)
	buffer.WriteString(" ")
	between.To.Write(buffer)
	buffer.WriteString(")")
}

// IsEmpty return true if statement is undefined.
func (between Between) IsEmpty() bool {
	return between.Identifier.IsEmpty() || between.Operator.IsEmpty() || between.And.IsEmpty() ||
		between.From == nil || between.To == nil || between.From.IsEmpty() || between.To.IsEmpty()
}
