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

	{
		query := loukoum.Update("table").Set(
			loukoum.Map{"a": 1, "b": 2},
			loukoum.Map{"c": 3, "d": 4})

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4", query.String())
	}
	{
		query := loukoum.Update("table").Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1})
		is.Equal("UPDATE table SET a = 1, foo = 2", query.String())
	}
}

func TestUpdate_Set_Map_Duplicates(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(
			loukoum.Map{"a": 1, "b": 2},
			loukoum.Map{"a": 3})

		is.Equal("UPDATE table SET a = 3, b = 2", query.String())
	}
	{
		query := loukoum.Update("table").Set(
			loukoum.Map{"a": 1, "b": 2},
			loukoum.Map{"b": 2})

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}
}

func TestUpdate_Set_Pair(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(loukoum.Pair("a", 1))
		is.Equal("UPDATE table SET a = 1", query.String())
	}
	{
		query := loukoum.Update("table").Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2))
		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}
	{
		query := loukoum.Update("table").Set(loukoum.Pair(loukoum.Column("a"), 1), loukoum.Pair("b", 2))
		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}
}

func TestUpdate_Set_Pair_Duplicates(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(
			loukoum.Pair("a", 1),
			loukoum.Pair("a", 1))

		is.Equal("UPDATE table SET a = 1", query.String())
	}
	{
		query := loukoum.Update("table").Set(
			loukoum.Pair("a", 1),
			loukoum.Pair("b", 2),
			loukoum.Pair("b", 2),
			loukoum.Pair("c", 3),
			loukoum.Pair("a", 4))

		is.Equal("UPDATE table SET a = 4, b = 2, c = 3", query.String())
	}
}

func TestUpdate_Set_MapAndPair(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(
			loukoum.Map{"a": 1, "b": 2},
			loukoum.Pair("d", 4),
			loukoum.Map{"c": 3},
			loukoum.Pair("e", 5))

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5", query.String())
	}
	{
		query := loukoum.Update("table").Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1})
		is.Equal("UPDATE table SET a = 1, foo = 2", query.String())
	}
}

func TestUpdate_Set_MapAndPair_Duplicates(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Update("table").Set(
			loukoum.Map{"a": 1, "b": 2},
			loukoum.Pair("a", 4),
			loukoum.Map{"b": 2},
			loukoum.Pair("c", 3))

		is.Equal("UPDATE table SET a = 4, b = 2, c = 3", query.String())
	}
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
