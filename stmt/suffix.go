package stmt

import (
	"github.com/ulule/loukoum/v2/types"
)

// Suffix is a suffix expression.
type Suffix struct {
	Suffix string
}

// NewSuffix returns a new Suffix instance.
func NewSuffix(suffix string) Suffix {
	return Suffix{
		Suffix: suffix,
	}
}

// Write exposes statement as a SQL query.
func (suffix Suffix) Write(ctx types.Context) {
	if suffix.IsEmpty() {
		return
	}
	ctx.Write(suffix.Suffix)
}

// IsEmpty returns true if statement is undefined.
func (suffix Suffix) IsEmpty() bool {
	return suffix.Suffix == ""
}

// Ensure that Suffix is a Statement
var _ Statement = Suffix{}
