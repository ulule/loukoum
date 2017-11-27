package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestUpdate(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table")
		is.Panics(func() { is.Equal("UPDATE table", query.String()) })
	}
	{
		query := loukoum.Update("table").Set(loukoum.Map{"b": "OK", "a": 1})
		is.Equal("UPDATE table SET a = 1, b = 'OK'", query.String())
	}
	{
		query := loukoum.Update("table").Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1})
		is.Equal("UPDATE table SET a = 1, foo = 2", query.String())
	}
}

func TestUpdate_Where(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(loukoum.Map{"a": 1}).Where(loukoum.Condition("id").Equal(1))
		is.Equal("UPDATE table SET a = 1 WHERE (id = 1)", query.String())
	}

	{
		query := loukoum.Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("status").Equal("online"))

		is.Equal("UPDATE table SET a = 1 WHERE ((id = 1) AND (status = 'online'))", query.String())
	}

	{
		query := loukoum.Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("status").Equal("online"))

		is.Equal("UPDATE table SET a = 1 WHERE ((id = 1) OR (status = 'online'))", query.String())
	}
}

func TestUpdate_From(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table1").
			Set(loukoum.Map{"a": 1}).
			From("table2").
			Where(loukoum.Condition("table2.id").Equal(loukoum.Raw("table1.id")))

		is.Equal("UPDATE table1 SET a = 1 FROM table2 WHERE (table2.id = table1.id)", query.String())
	}
}

func TestUpdate_Returning(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(loukoum.Map{"a": 1}).Returning("*")
		is.Equal("UPDATE table SET a = 1 RETURNING *", query.String())
	}
}
