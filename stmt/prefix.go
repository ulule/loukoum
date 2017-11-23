package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Prefix is a prefix expression.
type Prefix struct {
	Statement
	Prefix string
}

// NewPrefix returns a new Prefix instance.
func NewPrefix(prefix string) Prefix {
	return Prefix{
		Prefix: prefix,
	}
}

// Write exposes statement as a SQL query.
func (prefix Prefix) Write(ctx *types.Context) {
	if prefix.IsEmpty() {
		return
	}
	ctx.Write(prefix.Prefix)
}

// IsEmpty returns true if statement is undefined.
func (prefix Prefix) IsEmpty() bool {
	return prefix.Prefix == ""
}
