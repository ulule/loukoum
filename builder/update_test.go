package builder_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/format"
)

func TestUpdate_Set_Undefined(t *testing.T) {
	is := require.New(t)

	// Ensure that update panics if SET clause is empty.
	Failure(is, func() builder.Builder {
		return loukoum.Update("table")
	})
	Failure(is, func() builder.Builder {
		return loukoum.Update("table").Set("")
	})
}

func TestUpdate_Set_Map(t *testing.T) {
	is := require.New(t)

	// Variadic with built-in Map type
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Map{"c": 3, "d": 4},
			)

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4", query.String())
	}

	// Variadic with string / interface map
	{
		query := loukoum.
			Update("table").
			Set(
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"c": 3, "d": 4},
			)

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4", query.String())
	}

	// With column instance
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1})

		is.Equal("UPDATE table SET a = 1, foo = 2", query.String())
	}

	// Multiple Set() calls with built-in Map type
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Set(loukoum.Map{"b": 2}).
			Set(loukoum.Map{"c": 3}).
			Set(loukoum.Map{"d": 4, "e": 5})

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5", query.String())
	}

	// Multiple Set() calls with string / interface map.
	{
		query := loukoum.Update("table").
			Set(map[string]interface{}{"a": 1}).
			Set(map[string]interface{}{"b": 2}).
			Set(map[string]interface{}{"c": 3}).
			Set(map[string]interface{}{"d": 4, "e": 5})

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5", query.String())
	}

	// Multiple Set() calls with a mix of built-in Map type and string / interface map.
	{
		query := loukoum.Update("table").
			Set(loukoum.Map{"a": 1}).
			Set(map[string]interface{}{"b": 2}).
			Set(loukoum.Map{"c": 3}).
			Set(map[string]interface{}{"d": 4, "e": 5})

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5", query.String())
	}
}

func TestUpdate_Set_Map_Uniqueness(t *testing.T) {
	is := require.New(t)

	// Variadic with built-in Map type
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Map{"b": 2},
			)

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Variadic with string / interface map
	{
		query := loukoum.
			Update("table").
			Set(
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"b": 2},
			)

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Multiple Set() calls with built-in Map type
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(loukoum.Map{"b": 2})

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Multiple Set() calls with string / interface map
	{
		query := loukoum.
			Update("table").
			Set(map[string]interface{}{"a": 1, "b": 2}).
			Set(map[string]interface{}{"b": 2})

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Multiple Set() calls with a mix of built-in Map type and string / interface map.
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(map[string]interface{}{"b": 2})

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}
}

func TestUpdate_Set_Map_LastWriteWins(t *testing.T) {
	is := require.New(t)

	// Variadic with built-in Map type
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Map{"a": 3},
			)

		is.Equal("UPDATE table SET a = 3, b = 2", query.String())
	}

	// Variadic with string / interface map
	{
		query := loukoum.
			Update("table").
			Set(
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"a": 3},
			)

		is.Equal("UPDATE table SET a = 3, b = 2", query.String())
	}

	// Multiple Set() calls with built-in Map type
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(loukoum.Map{"a": 3})

		is.Equal("UPDATE table SET a = 3, b = 2", query.String())
	}

	// Multiple Set() calls with string / interface map
	{
		query := loukoum.
			Update("table").
			Set(map[string]interface{}{"a": 1, "b": 2}).
			Set(map[string]interface{}{"a": 3})

		is.Equal("UPDATE table SET a = 3, b = 2", query.String())
	}

	// Multiple Set() calls with a mix of built-in Map type and string / interface map.
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(map[string]interface{}{"a": 3})

		is.Equal("UPDATE table SET a = 3, b = 2", query.String())
	}
}

func TestUpdate_Set_Pair(t *testing.T) {
	is := require.New(t)

	// Single Pair
	{
		query := loukoum.Update("table").Set(loukoum.Pair("a", 1))
		is.Equal("UPDATE table SET a = 1", query.String())
	}

	// Variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("b", 2),
			)

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// Multiple Set() calls
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1)).
			Set(loukoum.Pair("b", 2))

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}

	// With column instance
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Pair(loukoum.Column("a"), 1),
				loukoum.Pair("b", 2),
			)

		is.Equal("UPDATE table SET a = 1, b = 2", query.String())
	}
}

func TestUpdate_Set_Pair_Uniqueness(t *testing.T) {
	is := require.New(t)

	// Variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("a", 1),
			)

		is.Equal("UPDATE table SET a = 1", query.String())
	}

	// Multiple Set() calls
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1), loukoum.Pair("a", 1)).
			Set(loukoum.Pair("a", 1), loukoum.Pair("a", 1))

		is.Equal("UPDATE table SET a = 1", query.String())
	}

	// Last write with variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("b", 2),
				loukoum.Pair("b", 5),
				loukoum.Pair("c", 3),
				loukoum.Pair("a", 4),
			)

		is.Equal("UPDATE table SET a = 4, b = 5, c = 3", query.String())
	}

	// Last write with multiple Set()
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2)).
			Set(loukoum.Pair("b", 5), loukoum.Pair("c", 3)).
			Set(loukoum.Pair("a", 4))

		is.Equal("UPDATE table SET a = 4, b = 5, c = 3", query.String())
	}
}

func TestUpdate_Set_Pair_LastWriteWins(t *testing.T) {
	is := require.New(t)

	// Variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("b", 2),
				loukoum.Pair("b", 5),
				loukoum.Pair("c", 3),
				loukoum.Pair("a", 4),
			)

		is.Equal("UPDATE table SET a = 4, b = 5, c = 3", query.String())
	}

	// Multiple Set() calls
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2)).
			Set(loukoum.Pair("b", 5), loukoum.Pair("c", 3)).
			Set(loukoum.Pair("a", 4))

		is.Equal("UPDATE table SET a = 4, b = 5, c = 3", query.String())
	}
}

func TestUpdate_Set_MapAndPair(t *testing.T) {
	is := require.New(t)

	// Variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Pair("d", 4),
				loukoum.Map{"c": 3},
				loukoum.Pair("e", 5),
			)

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5", query.String())
	}

	// Multiple Set() calls
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("d", 4)).
			Set(loukoum.Map{"c": 3}, loukoum.Pair("e", 5))

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5", query.String())
	}

	// With columns instances
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1})

		is.Equal("UPDATE table SET a = 1, foo = 2", query.String())
	}
}

func TestUpdate_Set_MapAndPair_Uniqueness(t *testing.T) {
	is := require.New(t)

	// Variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Pair("a", 1),
				loukoum.Map{"b": 2},
				loukoum.Pair("c", 3),
			)

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3", query.String())
	}

	// Multiple Set() calls
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("a", 1)).
			Set(loukoum.Map{"b": 2}, loukoum.Pair("b", 2)).
			Set(loukoum.Map{"c": 3}, loukoum.Pair("c", 3))

		is.Equal("UPDATE table SET a = 1, b = 2, c = 3", query.String())
	}
}

func TestUpdate_Set_MapAndPair_LastWriteWins(t *testing.T) {
	is := require.New(t)

	// Variadic
	{
		query := loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Pair("a", 2),
				loukoum.Map{"b": 3},
				loukoum.Pair("c", 4),
			)

		is.Equal("UPDATE table SET a = 2, b = 3, c = 4", query.String())
	}

	// Multiple Set() calls
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("a", 2)).
			Set(loukoum.Map{"b": 2}, loukoum.Pair("b", 3)).
			Set(loukoum.Map{"c": 3}, loukoum.Pair("c", 4))

		is.Equal("UPDATE table SET a = 2, b = 3, c = 4", query.String())
	}
}

func TestUpdate_Set_Valuer(t *testing.T) {
	is := require.New(t)

	// pq.NullTime
	{
		now := time.Now()
		fnow := format.Time(now)

		query := loukoum.Update("table").Set(loukoum.Map{"created_at": pq.NullTime{Time: now, Valid: true}})
		is.Equal(fmt.Sprintf("UPDATE table SET created_at = %s", fnow), query.String())

		query = loukoum.Update("table").Set(loukoum.Map{"created_at": pq.NullTime{Time: now, Valid: false}})
		is.Equal("UPDATE table SET created_at = NULL", query.String())
	}

	// sql.NullString
	{
		query := loukoum.Update("table").Set(loukoum.Map{"message": sql.NullString{String: "ok", Valid: true}})
		is.Equal("UPDATE table SET message = 'ok'", query.String())

		query = loukoum.Update("table").Set(loukoum.Map{"message": sql.NullString{String: "ok", Valid: false}})
		is.Equal("UPDATE table SET message = NULL", query.String())

	}

	// sql.NullInt
	{
		query := loukoum.Update("table").Set(loukoum.Map{"count": sql.NullInt64{Int64: 30, Valid: true}})
		is.Equal("UPDATE table SET count = 30", query.String())

		query = loukoum.Update("table").Set(loukoum.Map{"count": sql.NullInt64{Int64: 30, Valid: false}})
		is.Equal("UPDATE table SET count = NULL", query.String())
	}
}

func TestUpdate_Set_Using(t *testing.T) {
	is := require.New(t)

	// Invoking Using() without any argument must panic
	Failure(is, func() builder.Builder {
		return loukoum.Update("table").Set("a", "b", "c").Using()
	})

	// Invoking Using() without any columns must panic
	Failure(is, func() builder.Builder {
		return loukoum.Update("table").Set(loukoum.Pair("a", 30)).Using()
	})

	// Multi-values
	{
		query := loukoum.Update("table").Set("a", "b", "c").Using("d", "e", "f")
		is.Equal("UPDATE table SET (a, b, c) = ('d', 'e', 'f')", query.String())
	}

	// Columns to columns
	{
		query := loukoum.Update("table").Set("a", "b", "c").Using(loukoum.Raw("d+1"), loukoum.Raw("e+1"), loukoum.Raw("f+1"))
		is.Equal("UPDATE table SET (a, b, c) = (d+1, e+1, f+1)", query.String())
	}

	// Sub-select
	{
		query := loukoum.Update("table").
			Set("a", "b", "c").
			Using(
				loukoum.Select("a", "b", "c").
					From("table").
					Where(loukoum.Condition("disabled").Equal(false)).Statement(),
			)

		is.Equal("UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = false))", query.String())
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

	// Simple Where
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1))

		is.Equal("UPDATE table SET a = 1 WHERE (id = 1)", query.String())
	}

	// AND
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("status").Equal("online"))

		is.Equal("UPDATE table SET a = 1 WHERE ((id = 1) AND (status = 'online'))", query.String())
	}

	// OR
	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("status").Equal("online"))

		is.Equal("UPDATE table SET a = 1 WHERE ((id = 1) OR (status = 'online'))", query.String())
	}
}

func TestUpdate_From(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Update("table1").
			Set(loukoum.Map{"a": 1}).
			From("table2").
			Where(loukoum.Condition("table2.id").Equal(loukoum.Raw("table1.id")))

		is.Equal("UPDATE table1 SET a = 1 FROM table2 WHERE (table2.id = table1.id)", query.String())
	}
}

func TestUpdate_Returning(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Returning("*")

		is.Equal("UPDATE table SET a = 1 RETURNING *", query.String())
	}
}
