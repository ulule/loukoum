package builder_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
)

func TestUpdate_Set_Undefined(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "panic if set clause is empty 0",
			Failure: func() builder.Builder {
				return loukoum.Update("table")
			},
		},
		{
			Name: "panic if set clause is empty 1",
			Failure: func() builder.Builder {
				return loukoum.Update("table").Set("")
			},
		},
	})
}

func TestUpdate_Set_Map(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Set variadic",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Map{"a": 1, "b": 2},
						loukoum.Map{"c": 3, "d": 4},
					),
				loukoum.
					Update("table").
					Set(
						map[string]interface{}{"a": 1, "b": 2},
						map[string]interface{}{"c": 3, "d": 4},
					),
			},
			String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4",
			Args:       []interface{}{1, 2, 3, 4},
		},
		{
			Name: "With column instance",
			Builder: loukoum.
				Update("table").
				Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1}),
			String:     "UPDATE table SET a = 1, foo = 2",
			Query:      "UPDATE table SET a = $1, foo = $2",
			NamedQuery: "UPDATE table SET a = :arg_1, foo = :arg_2",
			Args:       []interface{}{1, 2},
		},
		{
			Name: "Multiple Set() calls with",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1}).
					Set(loukoum.Map{"b": 2}).
					Set(loukoum.Map{"c": 3}).
					Set(loukoum.Map{"d": 4, "e": 5}),
				loukoum.Update("table").
					Set(map[string]interface{}{"a": 1}).
					Set(map[string]interface{}{"b": 2}).
					Set(map[string]interface{}{"c": 3}).
					Set(map[string]interface{}{"d": 4, "e": 5}),
				loukoum.Update("table").
					Set(loukoum.Map{"a": 1}).
					Set(map[string]interface{}{"b": 2}).
					Set(loukoum.Map{"c": 3}).
					Set(map[string]interface{}{"d": 4, "e": 5}),
			},
			String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4, e = $5",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4, e = :arg_5",
			Args:       []interface{}{1, 2, 3, 4, 5},
		},
	})
}

func TestUpdate_Set_Map_Uniqueness(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Map{"a": 1, "b": 2},
						loukoum.Map{"b": 2},
					),
				loukoum.
					Update("table").
					Set(
						map[string]interface{}{"a": 1, "b": 2},
						map[string]interface{}{"b": 2},
					),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}).
					Set(loukoum.Map{"b": 2}),
				loukoum.
					Update("table").
					Set(map[string]interface{}{"a": 1, "b": 2}).
					Set(map[string]interface{}{"b": 2}),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}).
					Set(map[string]interface{}{"b": 2}),
			},
			String:     "UPDATE table SET a = 1, b = 2",
			Query:      "UPDATE table SET a = $1, b = $2",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
			Args:       []interface{}{1, 2},
		},
	})
}

func TestUpdate_Set_Map_LastWriteWins(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{

				loukoum.
					Update("table").
					Set(
						loukoum.Map{"a": 1, "b": 2},
						loukoum.Map{"a": 3},
					),
				loukoum.
					Update("table").
					Set(
						map[string]interface{}{"a": 1, "b": 2},
						map[string]interface{}{"a": 3},
					),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}).
					Set(loukoum.Map{"a": 3}),
				loukoum.
					Update("table").
					Set(map[string]interface{}{"a": 1, "b": 2}).
					Set(map[string]interface{}{"a": 3}),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}).
					Set(map[string]interface{}{"a": 3}),
			},
			String:     "UPDATE table SET a = 3, b = 2",
			Query:      "UPDATE table SET a = $1, b = $2",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
			Args:       []interface{}{3, 2},
		},
	})
}

func TestUpdate_Set_Pair(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:       "Single Pair",
			Builder:    loukoum.Update("table").Set(loukoum.Pair("a", 1)),
			String:     "UPDATE table SET a = 1",
			Query:      "UPDATE table SET a = $1",
			NamedQuery: "UPDATE table SET a = :arg_1",
			Args:       []interface{}{1},
		},
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Pair("a", 1),
						loukoum.Pair("b", 2),
					),
				loukoum.
					Update("table").
					Set(loukoum.Pair("a", 1)).
					Set(loukoum.Pair("b", 2)),
				loukoum.
					Update("table").
					Set(
						loukoum.Pair(loukoum.Column("a"), 1),
						loukoum.Pair("b", 2),
					),
			},
			String:     "UPDATE table SET a = 1, b = 2",
			Query:      "UPDATE table SET a = $1, b = $2",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
			Args:       []interface{}{1, 2},
		},
	})
}

func TestUpdate_Set_Pair_Uniqueness(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Pair("a", 1),
						loukoum.Pair("a", 1),
					),
				loukoum.
					Update("table").
					Set(loukoum.Pair("a", 1), loukoum.Pair("a", 1)).
					Set(loukoum.Pair("a", 1), loukoum.Pair("a", 1)),
			},
			String:     "UPDATE table SET a = 1",
			Query:      "UPDATE table SET a = $1",
			NamedQuery: "UPDATE table SET a = :arg_1",
			Args:       []interface{}{1},
		},
		{
			Name: "Last write",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Pair("a", 1),
						loukoum.Pair("b", 2),
						loukoum.Pair("b", 5),
						loukoum.Pair("c", 3),
						loukoum.Pair("a", 4),
					),
				loukoum.
					Update("table").
					Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2)).
					Set(loukoum.Pair("b", 5), loukoum.Pair("c", 3)).
					Set(loukoum.Pair("a", 4)),
			},
			String:     "UPDATE table SET a = 4, b = 5, c = 3",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
			Args:       []interface{}{4, 5, 3},
		},
	})
}

func TestUpdate_Set_Pair_LastWriteWins(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Pair("a", 1),
						loukoum.Pair("b", 2),
						loukoum.Pair("b", 5),
						loukoum.Pair("c", 3),
						loukoum.Pair("a", 4),
					),
				loukoum.
					Update("table").
					Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2)).
					Set(loukoum.Pair("b", 5), loukoum.Pair("c", 3)).
					Set(loukoum.Pair("a", 4)),
			},
			String:     "UPDATE table SET a = 4, b = 5, c = 3",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
			Args:       []interface{}{4, 5, 3},
		},
	})
}

func TestUpdate_Set_MapAndPair(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Map{"a": 1, "b": 2},
						loukoum.Pair("d", 4),
						loukoum.Map{"c": 3},
						loukoum.Pair("e", 5),
					),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("d", 4)).
					Set(loukoum.Map{"c": 3}, loukoum.Pair("e", 5)),
			},
			String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4, e = $5",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4, e = :arg_5",
			Args:       []interface{}{1, 2, 3, 4, 5},
		},
		{
			Name: "With column instances",
			Builder: loukoum.
				Update("table").
				Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1}),
			String:     "UPDATE table SET a = 1, foo = 2",
			Query:      "UPDATE table SET a = $1, foo = $2",
			NamedQuery: "UPDATE table SET a = :arg_1, foo = :arg_2",
			Args:       []interface{}{1, 2},
		},
	})
}

func TestUpdate_Set_MapAndPair_Uniqueness(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Map{"a": 1, "b": 2},
						loukoum.Pair("a", 1),
						loukoum.Map{"b": 2},
						loukoum.Pair("c", 3),
					),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("a", 1)).
					Set(loukoum.Map{"b": 2}, loukoum.Pair("b", 2)).
					Set(loukoum.Map{"c": 3}, loukoum.Pair("c", 3)),
			},
			String:     "UPDATE table SET a = 1, b = 2, c = 3",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
			Args:       []interface{}{1, 2, 3},
		},
	})
}

func TestUpdate_Set_MapAndPair_LastWriteWins(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set(
						loukoum.Map{"a": 1, "b": 2},
						loukoum.Pair("a", 2),
						loukoum.Map{"b": 3},
						loukoum.Pair("c", 4),
					),
				loukoum.
					Update("table").
					Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("a", 2)).
					Set(loukoum.Map{"b": 2}, loukoum.Pair("b", 3)).
					Set(loukoum.Map{"c": 3}, loukoum.Pair("c", 4)),
			},
			String:     "UPDATE table SET a = 2, b = 3, c = 4",
			Query:      "UPDATE table SET a = $1, b = $2, c = $3",
			NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
			Args:       []interface{}{2, 3, 4},
		},
	})
}

func TestUpdate_Set_Valuer(t *testing.T) {
	when, err := time.Parse(time.RFC3339, "2017-11-23T17:47:27+01:00")
	if err != nil {
		t.Fatal(err)
	}
	RunBuilderTests(t, []BuilderTest{
		{
			Name:       "pq.NullTime not null",
			Builder:    loukoum.Update("table").Set(loukoum.Map{"created_at": pq.NullTime{Time: when, Valid: true}}),
			String:     "UPDATE table SET created_at = '2017-11-23 16:47:27+00'",
			Query:      "UPDATE table SET created_at = $1",
			NamedQuery: "UPDATE table SET created_at = :arg_1",
			Args:       []interface{}{pq.NullTime{Time: when, Valid: true}},
		},
		{
			Name:       "pq.NullTime null",
			Builder:    loukoum.Update("table").Set(loukoum.Map{"created_at": pq.NullTime{Time: when, Valid: false}}),
			String:     "UPDATE table SET created_at = NULL",
			Query:      "UPDATE table SET created_at = $1",
			NamedQuery: "UPDATE table SET created_at = :arg_1",
			Args:       []interface{}{pq.NullTime{Time: when, Valid: false}},
		},
		{
			Name:       "pq.NullString not null",
			Builder:    loukoum.Update("table").Set(loukoum.Map{"message": sql.NullString{String: "ok", Valid: true}}),
			String:     "UPDATE table SET message = 'ok'",
			Query:      "UPDATE table SET message = $1",
			NamedQuery: "UPDATE table SET message = :arg_1",
			Args:       []interface{}{sql.NullString{String: "ok", Valid: true}},
		},
		{
			Name:       "pq.NullString null",
			Builder:    loukoum.Update("table").Set(loukoum.Map{"message": sql.NullString{String: "ok", Valid: false}}),
			String:     "UPDATE table SET message = NULL",
			Query:      "UPDATE table SET message = $1",
			NamedQuery: "UPDATE table SET message = :arg_1",
			Args:       []interface{}{sql.NullString{String: "ok", Valid: false}},
		},
		{
			Name:       "sql.NullInt64 not null",
			Builder:    loukoum.Update("table").Set(loukoum.Map{"count": sql.NullInt64{Int64: 30, Valid: true}}),
			String:     "UPDATE table SET count = 30",
			Query:      "UPDATE table SET count = $1",
			NamedQuery: "UPDATE table SET count = :arg_1",
			Args:       []interface{}{sql.NullInt64{Int64: 30, Valid: true}},
		},
		{
			Name:       "sql.NullInt64 null",
			Builder:    loukoum.Update("table").Set(loukoum.Map{"count": sql.NullInt64{Int64: 30, Valid: false}}),
			String:     "UPDATE table SET count = NULL",
			Query:      "UPDATE table SET count = $1",
			NamedQuery: "UPDATE table SET count = :arg_1",
			Args:       []interface{}{sql.NullInt64{Int64: 30, Valid: false}},
		},
	})
}

func TestUpdate_Set_Subquery(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Coalesce",
			Builder: loukoum.Update("test1").
				Set(loukoum.Pair("count",
					loukoum.Select("COALESCE(COUNT(*), 0)").
						From("test2").
						Where(loukoum.Condition("disabled").Equal(false)),
				)),
			String: fmt.Sprint(
				"UPDATE test1 SET count = (SELECT COALESCE(COUNT(*), 0)",
				" FROM test2 WHERE (disabled = false))",
			),
			Query: fmt.Sprint(
				"UPDATE test1 SET count = (SELECT COALESCE(COUNT(*), 0)",
				" FROM test2 WHERE (disabled = $1))",
			),
			NamedQuery: fmt.Sprint(
				"UPDATE test1 SET count = (SELECT COALESCE(COUNT(*), 0)",
				" FROM test2 WHERE (disabled = :arg_1))",
			),
			Args: []interface{}{false},
		},
	})
}

func TestUpdate_Set_Using(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Invoking Using() without any argument must panic",
			Failure: func() builder.Builder {
				return loukoum.Update("table").Set("a", "b", "c").Using()
			},
		},
		{
			Name: "Invoking Using() without any columns must panic",
			Failure: func() builder.Builder {
				return loukoum.Update("table").Set(loukoum.Pair("a", 30)).Using()
			},
		},
		{
			Name: "Multi-values",
			Builders: []builder.Builder{
				loukoum.
					Update("table").
					Set("a", "b", "c").
					Using("d", "e", "f"),
				loukoum.
					Update("table").
					Set("a").
					Set("b").
					Set("c").
					Using("d", "e", "f"),
			},
			String:     "UPDATE table SET (a, b, c) = ('d', 'e', 'f')",
			Query:      "UPDATE table SET (a, b, c) = ($1, $2, $3)",
			NamedQuery: "UPDATE table SET (a, b, c) = (:arg_1, :arg_2, :arg_3)",
			Args:       []interface{}{"d", "e", "f"},
		},
		{
			Name: "Columns to columns",
			Builder: loukoum.
				Update("table").
				Set("a", "b", "c").
				Using(loukoum.Raw("d+1"), loukoum.Raw("e+1"), loukoum.Raw("f+1")),
			SameQuery: "UPDATE table SET (a, b, c) = (d+1, e+1, f+1)",
		},
		{
			Name: "Sub-select",
			Builder: loukoum.Update("table").
				Set("a", "b", "c").
				Using(
					loukoum.Select("a", "b", "c").
						From("table").
						Where(loukoum.Condition("disabled").Equal(false)),
				),
			String:     "UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = false))",
			Query:      "UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = $1))",
			NamedQuery: "UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = :arg_1))",
			Args:       []interface{}{false},
		},
	})
}

func TestUpdate_OnlyTable(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:       "Simple",
			Builder:    loukoum.Update("table").Only().Set(loukoum.Map{"a": 1}),
			String:     "UPDATE ONLY table SET a = 1",
			Query:      "UPDATE ONLY table SET a = $1",
			NamedQuery: "UPDATE ONLY table SET a = :arg_1",
			Args:       []interface{}{1},
		},
	})
}

func TestUpdate_Where(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builder: loukoum.
				Update("table").
				Set(loukoum.Map{"a": 1}).
				Where(loukoum.Condition("id").Equal(1)),
			String:     "UPDATE table SET a = 1 WHERE (id = 1)",
			Query:      "UPDATE table SET a = $1 WHERE (id = $2)",
			NamedQuery: "UPDATE table SET a = :arg_1 WHERE (id = :arg_2)",
			Args:       []interface{}{1, 1},
		},
		{
			Name: "AND",
			Builder: loukoum.
				Update("table").
				Set(loukoum.Map{"a": 1}).
				Where(loukoum.Condition("id").Equal(1)).
				And(loukoum.Condition("status").Equal("online")),
			String:     "UPDATE table SET a = 1 WHERE ((id = 1) AND (status = 'online'))",
			Query:      "UPDATE table SET a = $1 WHERE ((id = $2) AND (status = $3))",
			NamedQuery: "UPDATE table SET a = :arg_1 WHERE ((id = :arg_2) AND (status = :arg_3))",
			Args:       []interface{}{1, 1, "online"},
		},
		{
			Name: "OR",
			Builder: loukoum.
				Update("table").
				Set(loukoum.Map{"a": 1}).
				Where(loukoum.Condition("id").Equal(1)).
				Or(loukoum.Condition("status").Equal("online")),
			String:     "UPDATE table SET a = 1 WHERE ((id = 1) OR (status = 'online'))",
			Query:      "UPDATE table SET a = $1 WHERE ((id = $2) OR (status = $3))",
			NamedQuery: "UPDATE table SET a = :arg_1 WHERE ((id = :arg_2) OR (status = :arg_3))",
			Args:       []interface{}{1, 1, "online"},
		},
	})
}

func TestUpdate_From(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builder: loukoum.
				Update("table1").
				Set(loukoum.Map{"a": 1}).
				From("table2").
				Where(loukoum.Condition("table2.id").Equal(loukoum.Raw("table1.id"))),
			String:     "UPDATE table1 SET a = 1 FROM table2 WHERE (table2.id = table1.id)",
			Query:      "UPDATE table1 SET a = $1 FROM table2 WHERE (table2.id = table1.id)",
			NamedQuery: "UPDATE table1 SET a = :arg_1 FROM table2 WHERE (table2.id = table1.id)",
			Args:       []interface{}{1},
		},
	})
}

func TestUpdate_Returning(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "*",
			Builder: loukoum.
				Update("table").
				Set(loukoum.Map{"a": 1}).
				Returning("*"),
			String:     "UPDATE table SET a = 1 RETURNING *",
			Query:      "UPDATE table SET a = $1 RETURNING *",
			NamedQuery: "UPDATE table SET a = :arg_1 RETURNING *",
			Args:       []interface{}{1},
		},
	})
}
