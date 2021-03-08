package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Values is a VALUES clause.
type Values struct {
	Values Expression
}

// NewValues returns a new Values instance.
func NewValues(values Expression) Values {
	return Values{
		Values: values,
	}
}

// Write exposes statement as a SQL query.
func (values Values) Write(ctx types.Context) {
	if values.IsEmpty() {
		return
	}

	ctx.Write(token.Values.String())
	ctx.Write(" (")
	values.Values.Write(ctx)
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (values Values) IsEmpty() bool {
	return values.Values == nil || (values.Values != nil && values.Values.IsEmpty())
}

// Ensure that Values is a Statement
var _ Statement = Values{}
