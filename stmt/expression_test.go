package stmt_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/v3/stmt"
	"github.com/ulule/loukoum/v3/types"
)

func TestExpression_Valuer(t *testing.T) {
	is := require.New(t)

	// pq.NullTime
	{
		ctx := &types.NamedContext{}

		source := pq.NullTime{Valid: true, Time: time.Now()}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source, args["arg_1"])
	}
	{
		ctx := &types.NamedContext{}

		source := pq.NullTime{}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source, args["arg_1"])
	}

	// sql.NullString
	{
		ctx := &types.NamedContext{}

		source := sql.NullString{Valid: true, String: "ok"}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source, args["arg_1"])
	}
	{
		ctx := &types.NamedContext{}

		source := sql.NullString{}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source, args["arg_1"])
	}

	// sql.NullInt64
	{
		ctx := &types.NamedContext{}

		source := sql.NullInt64{Valid: true, Int64: 32}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source, args["arg_1"])
	}
	{
		ctx := &types.NamedContext{}

		source := sql.NullInt64{}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source, args["arg_1"])
	}
}

func TestExpression_Encoder(t *testing.T) {
	is := require.New(t)

	// StringEncoder
	{
		ctx := &types.NamedContext{}

		source := tsencoder{value: "foobar"}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

	// Int64Encoder
	{
		ctx := &types.NamedContext{}

		source := tiencoder{value: 32}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

	// BoolEncoder
	{
		ctx := &types.NamedContext{}

		source := tbencoder{value: true}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(ctx)
		query := ctx.Query()
		args := ctx.Values()

		is.Equal(":arg_1", query)
		is.Equal(source.value, args["arg_1"])
	}

	// TimeEncoder
	{
		ctx := &types.NamedContext{}

		source := ttencoder{value: time.Now()}
		expression := stmt.NewExpression(source)
		value := expression.(stmt.Value)

		is.Equal(source.value, value.Value)

		value.Write(ctx)
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
