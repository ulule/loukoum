package stmt

import (
	"strconv"
	"strings"

	"github.com/ulule/loukoum/v3/types"
)

// Statement is the interface of the component which is the minimum unit constituting SQL.
// All types that implement this interface can be built as SQL.
type Statement interface {
	// IsEmpty returns true if statement is undefined.
	IsEmpty() bool
	// Write exposes statement as a SQL query.
	Write(ctx types.Context)
}

func quote(ident string) string {
	split := strings.Split(ident, ".")
	quoted := make([]string, 0, len(split))
	for i := range split {
		quoted = append(quoted, strconv.Quote(split[i]))
	}
	return strings.Join(quoted, ".")
}
