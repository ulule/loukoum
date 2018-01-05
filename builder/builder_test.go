package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

func Failure(is *require.Assertions, callback func() builder.Builder) {
	is.Panics(func() {
		_ = callback().String()
	})
}

func TestToColumns(t *testing.T) {
	is := require.New(t)

	// Simple cases...
	{
		args := []interface{}{"foo", "bar"}
		columns := builder.ToColumns(args)
		is.Len(columns, 2)
		is.Equal(stmt.NewColumn("foo"), columns[0])
		is.Equal(stmt.NewColumn("bar"), columns[1])
	}
	{
		args := []interface{}{stmt.NewColumn("foo"), "bar"}
		columns := builder.ToColumns(args)
		is.Len(columns, 2)
		is.Equal(stmt.NewColumn("foo"), columns[0])
		is.Equal(stmt.NewColumn("bar"), columns[1])
	}
	{
		args := []interface{}{stmt.NewColumn("foo"), stmt.NewColumn("bar")}
		columns := builder.ToColumns(args)
		is.Len(columns, 2)
		is.Equal(stmt.NewColumn("foo"), columns[0])
		is.Equal(stmt.NewColumn("bar"), columns[1])
	}
	{
		args := []interface{}{[]string{"foo", "bar"}}
		columns := builder.ToColumns(args)
		is.Len(columns, 2)
		is.Equal(stmt.NewColumn("foo"), columns[0])
		is.Equal(stmt.NewColumn("bar"), columns[1])
	}
	{
		args := []interface{}{[]stmt.Column{stmt.NewColumn("foo"), stmt.NewColumn("bar")}}
		columns := builder.ToColumns(args)
		is.Len(columns, 2)
		is.Equal(stmt.NewColumn("foo"), columns[0])
		is.Equal(stmt.NewColumn("bar"), columns[1])
	}
	{
		args := []interface{}{"foo, bar"}
		columns := builder.ToColumns(args)
		is.Len(columns, 2)
		is.Equal(stmt.NewColumn("foo"), columns[0])
		is.Equal(stmt.NewColumn("bar"), columns[1])
	}

	// Corner cases...
	{
		is.Panics(func() {
			args := []interface{}{""}
			builder.ToColumns(args)
		})
	}
	{
		is.Panics(func() {
			args := []interface{}{stmt.NewColumn("")}
			builder.ToColumns(args)
		})
	}
	{
		is.Panics(func() {
			args := []interface{}{[]string{""}}
			builder.ToColumns(args)
		})
	}
	{
		is.Panics(func() {
			args := []interface{}{[]stmt.Column{stmt.NewColumn("")}}
			builder.ToColumns(args)
		})
	}
	{
		is.Panics(func() {
			args := []interface{}{","}
			builder.ToColumns(args)
		})
	}
}
