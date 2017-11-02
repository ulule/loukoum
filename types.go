package loukoum

import (
	"github.com/ulule/loukoum/stmt"
)

// Column is a wrapper to create a new Column statement.
func Column(name string) stmt.Column {
	return stmt.NewColumn(name)
}
