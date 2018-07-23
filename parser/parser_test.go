package parser_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/parser"
	"github.com/ulule/loukoum/stmt"
	"github.com/ulule/loukoum/types"
)

func TestParse_Select(t *testing.T) {
	is := require.New(t)

	{
		input := "SELECT * FROM table"
		result, err := parser.Parse(input)
		is.NoError(err)
		is.NotNil(result)
		is.IsType(stmt.Select{}, result)
		query := (result).(stmt.Select)

		is.True(query.Prefix.IsEmpty())
		is.True(query.With.IsEmpty())
		is.False(query.Distinct)
		is.Len(query.Columns, 1)
		is.Equal("*", query.Columns[0].Name)
		is.Equal("", query.Columns[0].Alias)
		is.False(query.From.Only)
		is.Equal("table", query.From.Table.Name)
		is.Equal("", query.From.Table.Alias)
		is.Len(query.Joins, 0)
		is.True(query.Where.IsEmpty())
		is.True(query.GroupBy.IsEmpty())
		is.True(query.Having.IsEmpty())
		is.True(query.OrderBy.IsEmpty())
		is.True(query.Limit.IsEmpty())
		is.True(query.Offset.IsEmpty())
		is.True(query.Suffix.IsEmpty())

		ctx := &types.RawContext{}
		result.Write(ctx)
		is.Equal(input, ctx.Query())
	}
	{
		input := "SELECT a, b, c FROM table"
		result, err := parser.Parse(input)
		is.NoError(err)
		is.NotNil(result)
		is.IsType(stmt.Select{}, result)
		query := (result).(stmt.Select)

		is.True(query.Prefix.IsEmpty())
		is.True(query.With.IsEmpty())
		is.False(query.Distinct)
		is.Len(query.Columns, 3)
		is.Equal("a", query.Columns[0].Name)
		is.Equal("", query.Columns[0].Alias)
		is.Equal("b", query.Columns[1].Name)
		is.Equal("", query.Columns[1].Alias)
		is.Equal("c", query.Columns[2].Name)
		is.Equal("", query.Columns[2].Alias)
		is.False(query.From.Only)
		is.Equal("table", query.From.Table.Name)
		is.Equal("", query.From.Table.Alias)
		is.Len(query.Joins, 0)
		is.True(query.Where.IsEmpty())
		is.True(query.GroupBy.IsEmpty())
		is.True(query.Having.IsEmpty())
		is.True(query.OrderBy.IsEmpty())
		is.True(query.Limit.IsEmpty())
		is.True(query.Offset.IsEmpty())
		is.True(query.Suffix.IsEmpty())

		ctx := &types.RawContext{}
		result.Write(ctx)
		is.Equal(input, ctx.Query())
	}
	{
		input := "SELECT a AS fa, b AS fb, c AS fc FROM table"
		result, err := parser.Parse(input)
		is.NoError(err)
		is.NotNil(result)
		is.IsType(stmt.Select{}, result)
		query := (result).(stmt.Select)

		is.True(query.Prefix.IsEmpty())
		is.True(query.With.IsEmpty())
		is.False(query.Distinct)
		is.Len(query.Columns, 3)
		is.Equal("a", query.Columns[0].Name)
		is.Equal("fa", query.Columns[0].Alias)
		is.Equal("b", query.Columns[1].Name)
		is.Equal("fb", query.Columns[1].Alias)
		is.Equal("c", query.Columns[2].Name)
		is.Equal("fc", query.Columns[2].Alias)
		is.False(query.From.Only)
		is.Equal("table", query.From.Table.Name)
		is.Equal("", query.From.Table.Alias)
		is.Len(query.Joins, 0)
		is.True(query.Where.IsEmpty())
		is.True(query.GroupBy.IsEmpty())
		is.True(query.Having.IsEmpty())
		is.True(query.OrderBy.IsEmpty())
		is.True(query.Limit.IsEmpty())
		is.True(query.Offset.IsEmpty())
		is.True(query.Suffix.IsEmpty())

		ctx := &types.RawContext{}
		query.Write(ctx)
		fmt.Println(ctx.Query())
	}
	{
		result, err := parser.Parse("SELECT a, b, c, d, e, f FROM table AS foobar")
		is.NoError(err)
		is.NotNil(result)
		is.IsType(stmt.Select{}, result)
		query := (result).(stmt.Select)

		is.True(query.Prefix.IsEmpty())
		is.True(query.With.IsEmpty())
		is.False(query.Distinct)
		is.Len(query.Columns, 6)
		is.Equal("a", query.Columns[0].Name)
		is.Equal("", query.Columns[0].Alias)
		is.Equal("b", query.Columns[1].Name)
		is.Equal("", query.Columns[1].Alias)
		is.Equal("c", query.Columns[2].Name)
		is.Equal("", query.Columns[2].Alias)
		is.Equal("d", query.Columns[3].Name)
		is.Equal("", query.Columns[3].Alias)
		is.Equal("e", query.Columns[4].Name)
		is.Equal("", query.Columns[4].Alias)
		is.Equal("f", query.Columns[5].Name)
		is.Equal("", query.Columns[5].Alias)
		is.False(query.From.Only)
		is.Equal("table", query.From.Table.Name)
		is.Equal("foobar", query.From.Table.Alias)
		is.Len(query.Joins, 0)
		is.True(query.Where.IsEmpty())
		is.True(query.GroupBy.IsEmpty())
		is.True(query.Having.IsEmpty())
		is.True(query.OrderBy.IsEmpty())
		is.True(query.Limit.IsEmpty())
		is.True(query.Offset.IsEmpty())
		is.True(query.Suffix.IsEmpty())

		ctx := &types.RawContext{}
		result.Write(ctx)
		fmt.Println(ctx.Query())
	}
	{
		result, err := parser.Parse("SELECT DISTINCT * FROM table")
		is.NoError(err)
		is.NotNil(result)
		is.IsType(stmt.Select{}, result)
		query := (result).(stmt.Select)

		is.True(query.Prefix.IsEmpty())
		is.True(query.With.IsEmpty())
		is.True(query.Distinct)
		is.Len(query.Columns, 1)
		is.Equal("*", query.Columns[0].Name)
		is.Equal("", query.Columns[0].Alias)
		is.False(query.From.Only)
		is.Equal("table", query.From.Table.Name)
		is.Equal("", query.From.Table.Alias)
		is.Len(query.Joins, 0)
		is.Len(query.Joins, 0)
		is.True(query.Where.IsEmpty())
		is.True(query.GroupBy.IsEmpty())
		is.True(query.Having.IsEmpty())
		is.True(query.OrderBy.IsEmpty())
		is.True(query.Limit.IsEmpty())
		is.True(query.Offset.IsEmpty())
		is.True(query.Suffix.IsEmpty())

		ctx := &types.RawContext{}
		result.Write(ctx)
		fmt.Println(ctx.Query())
	}
	{
		result, err := parser.Parse("SELECT DISTINCT a, b, c FROM table")
		is.NoError(err)
		is.NotNil(result)
		is.IsType(stmt.Select{}, result)
		query := (result).(stmt.Select)

		is.True(query.Prefix.IsEmpty())
		is.True(query.With.IsEmpty())
		is.True(query.Distinct)
		is.Len(query.Columns, 3)
		is.Equal("a", query.Columns[0].Name)
		is.Equal("", query.Columns[0].Alias)
		is.Equal("b", query.Columns[1].Name)
		is.Equal("", query.Columns[1].Alias)
		is.Equal("c", query.Columns[2].Name)
		is.Equal("", query.Columns[2].Alias)
		is.False(query.From.Only)
		is.Equal("table", query.From.Table.Name)
		is.Equal("", query.From.Table.Alias)
		is.Len(query.Joins, 0)
		is.True(query.Where.IsEmpty())
		is.True(query.GroupBy.IsEmpty())
		is.True(query.Having.IsEmpty())
		is.True(query.OrderBy.IsEmpty())
		is.True(query.Limit.IsEmpty())
		is.True(query.Offset.IsEmpty())
		is.True(query.Suffix.IsEmpty())

		ctx := &types.RawContext{}
		result.Write(ctx)
		fmt.Println(ctx.Query())
	}
}
