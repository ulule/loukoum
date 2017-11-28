package builder_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestUpdate_Set_Undefined(t *testing.T) {
	is := require.New(t)
	is.Panics(func() { _ = loukoum.Update("table").String() })
}

func TestUpdate_Set_Map(t *testing.T) {
	is := require.New(t)

	m := loukoum.Map{"a": 1, "b": 2}
	p := loukoum.Pair("b", 2)

	// Simple map
	{
		query := loukoum.Update("table").Set(m)
		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Map with Column instance
	{
		query := loukoum.Update("table").Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1})
		is.Equal("UPDATE table SET a = 1, foo = 2", query.String())
	}

	// Two maps -> panic
	is.Panics(func() { loukoum.Update("table").Set(m, m) })

	// Map and pair -> panic
	is.Panics(func() { loukoum.Update("table").Set(m, p) })
}

func TestUpdate_Set_Pair(t *testing.T) {
	is := require.New(t)

	p1 := loukoum.Pair("a", 1)
	p2 := loukoum.Pair("b", 2)

	// Variadic pairs
	{
		query := loukoum.Update("table").Set(p1, p2)
		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Pairs with Column instance
	{
		query := loukoum.Update("table").Set(loukoum.Pair(loukoum.Column("a"), 1), p2)
		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Pair and map -> panic
	is.Panics(func() { loukoum.Update("table").Set(loukoum.Pair("b", 2), loukoum.Map{"a": 1}) })
}

func TestUpdate_OnlyTable(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Only().Set(loukoum.Map{"a": 1})
		is.Equal("UPDATE ONLY table SET a = 1", query.String())
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
