package builder_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestDelete(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Delete("table")
		is.Equal("DELETE FROM table", query.String())
	}
	{
		query := loukoum.Delete("table").Only()
		is.Equal("DELETE FROM ONLY table", query.String())
	}
	{
		query := loukoum.Delete(loukoum.Table("table"))
		is.Equal("DELETE FROM table", query.String())
	}
	{
		query := loukoum.Delete(loukoum.Table("table")).Only()
		is.Equal("DELETE FROM ONLY table", query.String())
	}
	{
		query := loukoum.Delete(loukoum.Table("table").As("foobar"))
		is.Equal("DELETE FROM table AS foobar", query.String())
	}
	{
		query := loukoum.Delete(loukoum.Table("table").As("foobar")).Only()
		is.Equal("DELETE FROM ONLY table AS foobar", query.String())
	}
}

func TestDelete_Using(t *testing.T) {
	is := require.New(t)

	// One table
	{
		query := loukoum.Delete("table").Using("foobar")
		is.Equal("DELETE FROM table USING foobar", query.String())
	}
	{
		query := loukoum.Delete("table").Using(loukoum.Table("foobar"))
		is.Equal("DELETE FROM table USING foobar", query.String())
	}
	{
		query := loukoum.Delete("table").Using(loukoum.Table("foobar").As("foo"))
		is.Equal("DELETE FROM table USING foobar AS foo", query.String())
	}

	// Two tables
	{
		query := loukoum.Delete("table").Using("foobar", "example")
		is.Equal("DELETE FROM table USING foobar, example", query.String())
	}
	{
		query := loukoum.Delete("table").Using(loukoum.Table("foobar"), "example")
		is.Equal("DELETE FROM table USING foobar, example", query.String())
	}
	{
		query := loukoum.Delete("table").Using(loukoum.Table("example"), loukoum.Table("foobar").As("foo"))
		is.Equal("DELETE FROM table USING example, foobar AS foo", query.String())
	}
}

func TestDelete_Where(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Delete("table").Where(loukoum.Condition("id").Equal(1))
		is.Equal("DELETE FROM table WHERE (id = 1)", query.String())
	}
	{
		when, err := time.Parse(time.RFC3339, "2017-11-23T17:47:27+01:00")
		is.NoError(err)
		is.NotZero(when)

		query := loukoum.Delete("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("created_at").GreaterThan(when))

		is.Equal("DELETE FROM table WHERE ((id = 1) AND (created_at > '2017-11-23 16:47:27+00'))", query.String())
	}
}

func TestDelete_Returning(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Delete("table").Returning("*")
		is.Equal("DELETE FROM table RETURNING *", query.String())
	}
	{
		query := loukoum.Delete("table").Returning("id")
		is.Equal("DELETE FROM table RETURNING id", query.String())
	}
}
