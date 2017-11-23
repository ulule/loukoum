package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Suffix is a suffix expression.
type Suffix struct {
	Statement
	Suffix string
}

// NewSuffix returns a new Suffix instance.
func NewSuffix(suffix string) Suffix {
	return Suffix{
		Suffix: suffix,
	}
}

// Write exposes statement as a SQL query.
func (suffix Suffix) Write(ctx *types.Context) {
	if suffix.IsEmpty() {
		return
	}
	ctx.Write(suffix.Suffix)
}

// IsEmpty returns true if statement is undefined.
func (suffix Suffix) IsEmpty() bool {
	return suffix.Suffix == ""
}
