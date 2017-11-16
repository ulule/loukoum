package stmt

import (
	"bytes"
)

// Query is a query.
type Query interface {
	String() string
}

// Statement is the interface of the component
// which is the minimum unit constituting SQL.
// All types that implement this interface
// can be built as SQL.
type Statement interface {
	// IsEmpty return true if statement is undefined.
	IsEmpty() bool
	// Write expose statement as a SQL query.
	Write(buffer *bytes.Buffer)
}
