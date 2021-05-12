package parser_test

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/v3/parser"
	"github.com/ulule/loukoum/v3/stmt"
	"github.com/ulule/loukoum/v3/types"
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
		on, ok := query.Condition.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.user_id", on.Right.Name)
		is.Empty(on.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("INNER JOIN account ON (project.account_id = account.id)")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Equal("account", query.Table.Name)
		is.Empty(query.Table.Alias)
		on, ok := query.Condition.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("project.account_id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("account.id", on.Right.Name)
		is.Empty(on.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("RIGHT JOIN foobar ON foobar.group_id = test.group_id;")
		is.NoError(err)
		is.Equal(types.RightJoin, query.Type)
		is.Equal("foobar", query.Table.Name)
		is.Empty(query.Table.Alias)
		on, ok := query.Condition.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("foobar.group_id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("test.group_id", on.Right.Name)
		is.Empty(on.Right.Alias)
	}

	// Partials
	{
		query, err := parser.ParseJoin("ON user.id = project.user_id")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		on, ok := query.Condition.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.user_id", on.Right.Name)
		is.Empty(on.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("user.id = project.user_id")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		on, ok := query.Condition.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.user_id", on.Right.Name)
		is.Empty(on.Right.Alias)
	}

	// With logical operator
	{
		query, err := parser.ParseJoin("ON user.id = project.user_id AND user.hash = project.hash")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		infix, ok := query.Condition.(stmt.InfixExpression)
		is.True(ok)
		is.NotEmpty(infix)
		on, ok := infix.Left.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.user_id", on.Right.Name)
		is.Empty(on.Right.Alias)
		is.Equal(stmt.LogicalOperator{Operator: types.And}, infix.Operator)
		on, ok = infix.Right.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.hash", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.hash", on.Right.Name)
		is.Empty(on.Right.Alias)
	}
	{
		query, err := parser.ParseJoin("ON user.id = project.user_id OR user.hash = project.hash")
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		infix, ok := query.Condition.(stmt.InfixExpression)
		is.True(ok)
		is.NotEmpty(infix)
		on, ok := infix.Left.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.user_id", on.Right.Name)
		is.Empty(on.Right.Alias)
		is.Equal(stmt.LogicalOperator{Operator: types.Or}, infix.Operator)
		on, ok = infix.Right.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.hash", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.hash", on.Right.Name)
		is.Empty(on.Right.Alias)
	}
	{
		query, err := parser.ParseJoin(
			"ON user.id = project.user_id AND user.hash = project.hash OR user.group_id = project.group_id",
		)
		is.NoError(err)
		is.Equal(types.InnerJoin, query.Type)
		is.Empty(query.Table.Name)
		is.Empty(query.Table.Alias)
		infix1, ok := query.Condition.(stmt.InfixExpression)
		is.True(ok)
		is.NotEmpty(infix1)
		infix2, ok := infix1.Left.(stmt.InfixExpression)
		is.True(ok)
		is.NotEmpty(infix2)
		on, ok := infix2.Left.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.user_id", on.Right.Name)
		is.Empty(on.Right.Alias)
		is.Equal(stmt.LogicalOperator{Operator: types.And}, infix2.Operator)
		on, ok = infix2.Right.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.hash", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.hash", on.Right.Name)
		is.Empty(on.Right.Alias)
		is.Equal(stmt.LogicalOperator{Operator: types.Or}, infix1.Operator)
		on, ok = infix1.Right.(stmt.OnClause)
		is.True(ok)
		is.NotEmpty(on)
		is.Equal("user.group_id", on.Left.Name)
		is.Empty(on.Left.Alias)
		is.Equal("project.group_id", on.Right.Name)
		is.Empty(on.Right.Alias)
	}

	// Invalid
	{
		query, err := parser.ParseJoin("INNER JOIN account ON (project.account_id = *)")
		is.Error(err)
		is.Equal(parser.ErrJoinInvalidCondition, errors.Cause(err))
		is.Zero(query)
	}
}
