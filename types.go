package loukoum

import (
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

const (
	// InnerJoin is used for "INNER JOIN" in join statement.
	InnerJoin = types.InnerJoin
	// LeftJoin is used for "LEFT JOIN" in join statement.
	LeftJoin = types.LeftJoin
	// RightJoin is used for "RIGHT JOIN" in join statement.
	RightJoin = types.RightJoin
)

// Column is a wrapper to create a new Column statement.
func Column(name string) stmt.Column {
	return stmt.NewColumn(name)
}

// Table is a wrapper to create a new Table statement.
func Table(name string) stmt.Table {
	return stmt.NewTable(name)
}

// On is a wrapper to create a new On statement.
func On(left string, right string) stmt.On {
	return stmt.NewOn(stmt.NewColumn(left), stmt.NewColumn(right))
}
