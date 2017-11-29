package stmt_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

func TestExpression_Valuer(t *testing.T) {
	is := require.New(t)

	{
		s := sql.NullString{Valid: true, String: "ok"}
		e := stmt.NewExpression(s)
		v := e.(stmt.Value)
		is.Equal(s.String, v.Value)
	}
	{
		ctx := types.NewContext()

		s := sql.NullString{Valid: false, String: ""}
		e := stmt.NewExpression(s)

		v := e.(stmt.Value)
		v.Write(ctx)

		is.Equal("NULL", ctx.Query())
	}
}
