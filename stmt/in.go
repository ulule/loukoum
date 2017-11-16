package stmt

import (
	"bytes"

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

func (in In) Write(buffer *bytes.Buffer) {
	if in.IsEmpty() {
		panic("loukoum: expression is undefined")
	}

	buffer.WriteString("(")
	in.Identifier.Write(buffer)
	buffer.WriteString(" ")
	in.Operator.Write(buffer)
	buffer.WriteString(" ")
	in.Value.Write(buffer)
	buffer.WriteString(")")
}

// IsEmpty return true if statement is undefined.
func (in In) IsEmpty() bool {
	return in.Identifier.IsEmpty() || in.Operator.IsEmpty() || in.Value == nil || in.Value.IsEmpty()
}
