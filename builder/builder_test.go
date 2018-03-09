package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

type BuilderTest struct {
	Name       string
	Builder    builder.Builder
	String     string
	Query      string
	Args       []interface{}
	NamedQuery string
	NamedArgs  map[string]interface{}
	Failure    func() builder.Builder
}

func (b *BuilderTest) Run(t *testing.T) {
	if b.Failure != nil {
		t.Run("Failure", func(t *testing.T) {
			require.Panics(t, func() {
				_ = b.Failure()
			})
		})
	}
	if b.Builder == nil {
		return
	}
	t.Run("String", func(t *testing.T) {
		if b.String != "" {
			require.Equal(t, b.String, b.Builder.String())
		}
	})
	t.Run("Query", func(t *testing.T) {
		query, args := b.Builder.Query()
		if b.Query != "" {
			require.Equal(t, b.Query, query)
			require.Equal(t, b.Args, args)
		}
	})
	t.Run("NamedQuery", func(t *testing.T) {
		query, args := b.Builder.NamedQuery()
		if b.NamedQuery != "" {
			require.Equal(t, b.NamedQuery, query)
			require.Equal(t, b.NamedArgs, args)
		}
	})
}

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
