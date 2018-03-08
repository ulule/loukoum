package stmt

import (
	"github.com/ulule/loukoum/types"
)

// Statement is the interface of the component which is the minimum unit constituting SQL.
// All types that implement this interface can be built as SQL.
type Statement interface {
	// IsEmpty returns true if statement is undefined.
	IsEmpty() bool
	// Write exposes statement as a SQL query.
	Write(ctx types.Context)
}
