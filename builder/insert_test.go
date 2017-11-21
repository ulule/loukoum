package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/stmt"
)

func TestInsert_Columns(t *testing.T) {
	is := require.New(t)

	// With columns
	{
		query := loukoum.Insert("table").Columns("a", "b", "c")
		is.Equal("INSERT INTO table (a, b, c)", query.String())
	}

	// Without columns
	{
		query := loukoum.Insert("table")
		is.Equal("INSERT INTO table", query.String())
	}
}

func TestInsert_Values(t *testing.T) {
	is := require.New(t)

	// With columns
	{
		query := loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"})

		is.Equal("INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc')", query.String())
	}

	// Without columns
	{
		query := loukoum.
			Insert("table").
			Values([]string{"va", "vb", "vc"})

		is.Equal("INSERT INTO table VALUES ('va', 'vb', 'vc')", query.String())
	}
}

func TestInsert_Returning(t *testing.T) {
	is := require.New(t)

	// One column
	{
		query := loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a")

		is.Equal("INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a", query.String())
	}

	// Many columns
	{
		query := loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a", "b")

		is.Equal("INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING (a, b)", query.String())
	}

	// AS
	{
		query := loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning(stmt.NewColumnAlias("a", "alias_a"))

		is.Equal("INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a AS alias_a", query.String())
	}

	// TODO: expression
}
