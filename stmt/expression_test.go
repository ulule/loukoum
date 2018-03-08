package stmt_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

func TestExpression_Valuer(t *testing.T) {
	is := require.New(t)

	// pq.NullTime
	{
		var ctx types.NamedContext

		source := pq.NullTime{Valid: true, Time: time.Now()}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.Time, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.Time, args["arg_1"])
	}
	{
		var ctx types.NamedContext

		source := pq.NullTime{}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(nil, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal("NULL", query)
		is.Empty(args)
	}

	// sql.NullString
	{
		var ctx types.NamedContext

		source := sql.NullString{Valid: true, String: "ok"}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.String, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.String, args["arg_1"])
	}
	{
		var ctx types.NamedContext

		source := sql.NullString{}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(nil, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal("NULL", query)
		is.Empty(args)
	}

	// sql.NullInt64
	{
		var ctx types.NamedContext

		source := sql.NullInt64{Valid: true, Int64: 32}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.Int64, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.Int64, args["arg_1"])
	}
	{
		var ctx types.NamedContext

		source := sql.NullInt64{}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(nil, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal("NULL", query)
		is.Empty(args)
	}
}

func TestExpression_Encoder(t *testing.T) {
	is := require.New(t)

	// StringEncoder
	{
		var ctx types.NamedContext

		source := tsencoder{value: "foobar"}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

	// Int64Encoder
	{
		var ctx types.NamedContext

		source := tiencoder{value: 32}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

	// BoolEncoder
	{
		var ctx types.NamedContext

		source := tbencoder{value: true}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

	// TimeEncoder
	{
		var ctx types.NamedContext

		source := ttencoder{value: time.Now()}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(&ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

}

type tsencoder struct {
	value string
}

func (e tsencoder) String() string {
	return e.value
}

type tiencoder struct {
	value int64
}

func (e tiencoder) Int64() int64 {
	return e.value
}

type tbencoder struct {
	value bool
}

func (e tbencoder) Bool() bool {
	return e.value
}

type ttencoder struct {
	value time.Time
}

func (e ttencoder) Time() time.Time {
	return e.value
}
