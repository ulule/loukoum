package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// Values is the VALUES clause.
type Values struct {
	Statement
	Values Expression
}

// NewValues returns a new Values instance.
func NewValues(values Expression) Values {
	return Values{
		Values: values,
	}
}

// Write expose statement as a SQL query.
func (values Values) Write(ctx *types.Context) {
	if values.IsEmpty() {
		return
	}

	ctx.Write(token.Values.String())
	ctx.Write(" (")
	values.Values.Write(ctx)
	ctx.Write(")")
}

// IsEmpty return true if statement is undefined.
func (values Values) IsEmpty() bool {
	return values.Values == nil || (values.Values != nil && values.Values.IsEmpty())
}
