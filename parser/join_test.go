package parser_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/parser"
	"github.com/ulule/loukoum/types"
)

func TestParseJoin(t *testing.T) {
	is := require.New(t)

	// Expressions
	{
		query, err := parser.ParseJoin("LEFT JOIN project ON user.id = project.user_id")
		is.NoError(err)
		is.Equal(types.LeftJoin, query.Type)
		is.Equal("project", query.Table.Name)
		is.Empty(query.Table.Alias)
		is.Equal("user.id", query.Condition.Left.Name)
		is.Empty(query.Condition.Left.Alias)
		is.Equal("project.user_id", query.Condition.Right.Name)
		is.Empty(query.Condition.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("INNER JOIN account ON (project.account_id = account.id)")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Equal("account", query.Table.Name)
		is.Empty(query.Table.Alias)
		is.Equal("project.account_id", query.Condition.Left.Name)
		is.Empty(query.Condition.Left.Alias)
		is.Equal("account.id", query.Condition.Right.Name)
		is.Empty(query.Condition.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("RIGHT JOIN foobar ON foobar.group_id = test.group_id;")
		is.NoError(err)
		is.Equal(types.RightJoin, query.Type)
		is.Equal("foobar", query.Table.Name)
		is.Empty(query.Table.Alias)
		is.Equal("foobar.group_id", query.Condition.Left.Name)
		is.Empty(query.Condition.Left.Alias)
		is.Equal("test.group_id", query.Condition.Right.Name)
		is.Empty(query.Condition.Right.Alias)
	}

	// Partials
	{
		query, err := parser.ParseJoin("ON user.id = project.user_id")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		is.Equal("user.id", query.Condition.Left.Name)
		is.Empty(query.Condition.Left.Alias)
		is.Equal("project.user_id", query.Condition.Right.Name)
		is.Empty(query.Condition.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("user.id = project.user_id")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		is.Equal("user.id", query.Condition.Left.Name)
		is.Empty(query.Condition.Left.Alias)
		is.Equal("project.user_id", query.Condition.Right.Name)
		is.Empty(query.Condition.Right.Alias)
	}

	// Invalid
	{
		query, err := parser.ParseJoin("INNER JOIN account ON (project.account_id = *)")
		is.Error(err)
		is.Equal(parser.ErrJoinInvalidCondition, errors.Cause(err))
		is.Zero(query)
	}
}
