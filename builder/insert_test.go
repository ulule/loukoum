package builder_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
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

func TestInsert_OnConflict(t *testing.T) {
	is := require.New(t)

	// Do nothing without target
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict(loukoum.DoNothing())

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT DO NOTHING"), query.String())
	}

	// Do nothing
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", loukoum.DoNothing())

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email) DO NOTHING"), query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict(loukoum.Column("email"), loukoum.DoNothing())

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email) DO NOTHING"), query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", "uuid", loukoum.DoNothing())

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email, uuid) DO NOTHING"), query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", loukoum.Column("uuid"), "reference", loukoum.DoNothing())

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email, uuid, reference) DO NOTHING"), query.String())
	}

	// Do update
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", loukoum.DoUpdate(
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				loukoum.Pair("enabled", true),
			))

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email) DO UPDATE SET created_at = NOW(), enabled = true"), query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict(loukoum.Column("email"), loukoum.DoUpdate(
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				loukoum.Pair("enabled", true),
			))

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email) DO UPDATE SET created_at = NOW(), enabled = true"), query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", "uuid", loukoum.DoUpdate(
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				loukoum.Pair("enabled", true),
			))

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email, uuid) DO UPDATE SET created_at = NOW(), enabled = true"), query.String())
	}
	{
		query := loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", loukoum.Column("uuid"), "reference", loukoum.DoUpdate(
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				loukoum.Pair("enabled", true),
			))

		is.Equal(fmt.Sprint("INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email, uuid, reference) DO UPDATE SET created_at = NOW(), enabled = true"), query.String())
	}

	// Corner cases...
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict()
		})
	}
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email")
		})
	}
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict(loukoum.DoUpdate(
					loukoum.Pair("created_at", loukoum.Raw("NOW()")),
					loukoum.Pair("enabled", true),
				))
		})
	}
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", 6700)
		})
	}
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict(569)
		})
	}
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", "uuid")
		})
	}
	{
		Failure(is, func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict(loukoum.Column("email"), loukoum.Column("uuid"), loukoum.Column("reference"))
		})
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
