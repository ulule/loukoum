package builder_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

type BuilderTest struct {
	Name       string
	Builder    builder.Builder
	Builders   []builder.Builder
	SameQuery  string
	String     string
	Query      string
	NamedQuery string
	Args       []interface{}
	Failure    func() builder.Builder
}

func (b BuilderTest) builders() []builder.Builder {
	var builders []builder.Builder
	builders = append(builders, b.Builders...)
	if b.Builder != nil {
		builders = append(builders, b.Builder)
	}
	return builders
}

func toNamedArgs(args []interface{}) map[string]interface{} {
	if args == nil {
		return nil
	}
	named := make(map[string]interface{})
	for i, arg := range args {
		name := fmt.Sprintf("arg_%d", i+1)
		named[name] = arg
	}
	return named
}

func RunBuilderTests(t *testing.T, tests []BuilderTest) {
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Failure != nil {
				t.Run("Failure", func(t *testing.T) {
					require.Panics(t, func() {
						_ = tt.Failure().String()
					})
				})
				return
			}
			for i, builder := range tt.builders() {
				t.Run(strconv.Itoa(i), func(t *testing.T) {
					if tt.SameQuery != "" {
						tt.String = tt.SameQuery
						tt.Query = tt.SameQuery
						tt.NamedQuery = tt.SameQuery
					}
					t.Run("String", func(t *testing.T) {
						require.Equal(t, tt.String, builder.String())
					})
					t.Run("Query", func(t *testing.T) {
						query, args := builder.Query()
						require.Equal(t, tt.Query, query)
						require.Equal(t, tt.Args, args)
					})
					t.Run("NamedQuery", func(t *testing.T) {
						query, args := builder.NamedQuery()
						require.Equal(t, tt.NamedQuery, query)
						require.Equal(t, toNamedArgs(tt.Args), args)
					})
				})
			}
		})
	}
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
