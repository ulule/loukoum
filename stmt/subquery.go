package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Exists is a subquery expression.
type Exists struct {
	subquery Expression
}

// NewExists returns a new Exists instance.
func NewExists(value interface{}) Exists {
	return Exists{
		subquery: NewExpression(value),
	}
}

func (Exists) expression() {}

// Write exposes statement as a SQL query.
func (exists Exists) Write(ctx types.Context) {
	ctx.Write("(EXISTS (")
	exists.subquery.Write(ctx)
	ctx.Write("))")
}

// IsEmpty returns true if statement is undefined.
func (exists Exists) IsEmpty() bool {
	return false
}

// Ensure that Exists is an Expression
var _ Expression = Exists{}
