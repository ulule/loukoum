package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
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
	{
		query := loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values("va", "vb", "vc")

		is.Equal("INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc')", query.String())
	}

	// Without columns
	{
		query := loukoum.
			Insert("table").
			Values([]string{"va", "vb", "vc"})

		is.Equal("INSERT INTO table VALUES ('va', 'vb', 'vc')", query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Values("va", "vb", "vc")

		is.Equal("INSERT INTO table VALUES ('va', 'vb', 'vc')", query.String())
	}

	// With raw values
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()"))

		is.Equal("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW())", query.String())
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
			Returning(loukoum.Column("a").As("alias_a"))

		is.Equal("INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a AS alias_a", query.String())
	}

	// TODO: expression
}
