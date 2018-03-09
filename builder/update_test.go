package builder_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/lib/pq"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/format"
)

var updatetests = []BuilderTest{
	{
		Name: "Not Set",
		Failure: func() builder.Builder {
			return loukoum.Update("table")
		},
	},
	{
		Name: "Set empty string",
		Failure: func() builder.Builder {
			return loukoum.Update("table").Set("")
		},
	},
	{
		Name: "Set Map variadic with built-in Map type",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Map{"c": 3, "d": 4},
			),
		String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4",
		Args:       []interface{}{1, 2, 3, 4},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
			"arg_4": 4,
		},
	},
	{
		Name: "Set Map variadic with string / interface map",
		Builder: loukoum.
			Update("table").
			Set(
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"c": 3, "d": 4},
			),
		String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4",
		Args:       []interface{}{1, 2, 3, 4},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
			"arg_4": 4,
		},
	},
	{
		Name: "Set Map with column statement",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1}),
		String:     "UPDATE table SET a = 1, foo = 2",
		Query:      "UPDATE table SET a = $1, foo = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, foo = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map with multiple Set calls with build in Map type",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Set(loukoum.Map{"b": 2}).
			Set(loukoum.Map{"c": 3}).
			Set(loukoum.Map{"d": 4, "e": 5}),
		String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4, e = $5",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4, e = :arg_5",
		Args:       []interface{}{1, 2, 3, 4, 5},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
			"arg_4": 4,
			"arg_5": 5,
		},
	},
	{
		Name: "Set Map with multiple Set calls with string / interface map",
		Builder: loukoum.Update("table").
			Set(map[string]interface{}{"a": 1}).
			Set(map[string]interface{}{"b": 2}).
			Set(map[string]interface{}{"c": 3}).
			Set(map[string]interface{}{"d": 4, "e": 5}),
		String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4, e = $5",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4, e = :arg_5",
		Args:       []interface{}{1, 2, 3, 4, 5},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
			"arg_4": 4,
			"arg_5": 5,
		},
	},
	{
		Name: "Set Map with multipel Set calls with a mix of types",
		Builder: loukoum.Update("table").
			Set(loukoum.Map{"a": 1}).
			Set(map[string]interface{}{"b": 2}).
			Set(loukoum.Map{"c": 3}).
			Set(map[string]interface{}{"d": 4, "e": 5}),
		String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4, e = $5",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4, e = :arg_5",
		Args:       []interface{}{1, 2, 3, 4, 5},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
			"arg_4": 4,
			"arg_5": 5,
		},
	},
	{
		Name: "Set Map Uniqueness Variadic with build-in Map type",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Map{"b": 2},
			),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Uniqueness Variadic with string / interface map",
		Builder: loukoum.
			Update("table").
			Set(
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"b": 2},
			),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Uniqueness with Multiple Set() calls with built-in Map type",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(loukoum.Map{"b": 2}),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Uniqueness with multiple Set() calls with string / interface map",
		Builder: loukoum.
			Update("table").
			Set(map[string]interface{}{"a": 1, "b": 2}).
			Set(map[string]interface{}{"b": 2}),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Uniqueness multiple Set() calls with a mix of built-in Map type and string / interface map.",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(map[string]interface{}{"b": 2}),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Last Write Wins Variadic with built-in Map type",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Map{"a": 3},
			),
		String:     "UPDATE table SET a = 3, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{3, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 3,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Last Write Wins Variadic with string / interface map",
		Builder: loukoum.
			Update("table").
			Set(
				map[string]interface{}{"a": 1, "b": 2},
				map[string]interface{}{"a": 3},
			),
		String:     "UPDATE table SET a = 3, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{3, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 3,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Last Write Wins Multiple Set() calls with built-in Map type",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(loukoum.Map{"a": 3}),
		String:     "UPDATE table SET a = 3, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{3, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 3,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Last Write Wins Multiple Set() calls with string / interface map",
		Builder: loukoum.
			Update("table").
			Set(map[string]interface{}{"a": 1, "b": 2}).
			Set(map[string]interface{}{"a": 3}),
		String:     "UPDATE table SET a = 3, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{3, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 3,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map Last Write Wins Multiple Set() calls with a mix of built-in Map type and string / interface map.",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}).
			Set(map[string]interface{}{"a": 3}),
		String:     "UPDATE table SET a = 3, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{3, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 3,
			"arg_2": 2,
		},
	},
	{
		Name:       "Set Pair single",
		Builder:    loukoum.Update("table").Set(loukoum.Pair("a", 1)),
		String:     "UPDATE table SET a = 1",
		Query:      "UPDATE table SET a = $1",
		NamedQuery: "UPDATE table SET a = :arg_1",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Set Pair Variadic",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("b", 2),
			),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Pair Multiple Set() calls",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1)).
			Set(loukoum.Pair("b", 2)),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Pair With column instance",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Pair(loukoum.Column("a"), 1),
				loukoum.Pair("b", 2),
			),
		String:     "UPDATE table SET a = 1, b = 2",
		Query:      "UPDATE table SET a = $1, b = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Pair Uniqueness Variadic",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("a", 1),
			),
		String:     "UPDATE table SET a = 1",
		Query:      "UPDATE table SET a = $1",
		NamedQuery: "UPDATE table SET a = :arg_1",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Set Pair Uniqueness Multiple Set() calls",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1), loukoum.Pair("a", 1)).
			Set(loukoum.Pair("a", 1), loukoum.Pair("a", 1)),
		String:     "UPDATE table SET a = 1",
		Query:      "UPDATE table SET a = $1",
		NamedQuery: "UPDATE table SET a = :arg_1",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Set Pair Uniqueness Last write with variadic",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("b", 2),
				loukoum.Pair("b", 5),
				loukoum.Pair("c", 3),
				loukoum.Pair("a", 4),
			),
		String:     "UPDATE table SET a = 4, b = 5, c = 3",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{4, 5, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 4,
			"arg_2": 5,
			"arg_3": 3,
		},
	},
	{
		Name: "Set Pair Uniqueness Last write with multiple Set()",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2)).
			Set(loukoum.Pair("b", 5), loukoum.Pair("c", 3)).
			Set(loukoum.Pair("a", 4)),
		String:     "UPDATE table SET a = 4, b = 5, c = 3",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{4, 5, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 4,
			"arg_2": 5,
			"arg_3": 3,
		},
	},
	{
		Name: "Set Pair Last Write Wins Variadic",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Pair("a", 1),
				loukoum.Pair("b", 2),
				loukoum.Pair("b", 5),
				loukoum.Pair("c", 3),
				loukoum.Pair("a", 4),
			),
		String:     "UPDATE table SET a = 4, b = 5, c = 3",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{4, 5, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 4,
			"arg_2": 5,
			"arg_3": 3,
		},
	},
	{
		Name: "Set Pair Last Write Wins Multiple Set() calls",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Pair("a", 1), loukoum.Pair("b", 2)).
			Set(loukoum.Pair("b", 5), loukoum.Pair("c", 3)).
			Set(loukoum.Pair("a", 4)),
		String:     "UPDATE table SET a = 4, b = 5, c = 3",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{4, 5, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 4,
			"arg_2": 5,
			"arg_3": 3,
		},
	},
	{
		Name: "Set Map and Pair Multiple Set() calls",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("d", 4)).
			Set(loukoum.Map{"c": 3}, loukoum.Pair("e", 5)),
		String:     "UPDATE table SET a = 1, b = 2, c = 3, d = 4, e = 5",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3, d = $4, e = $5",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3, d = :arg_4, e = :arg_5",
		Args:       []interface{}{1, 2, 3, 4, 5},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
			"arg_4": 4,
			"arg_5": 5,
		},
	},
	{
		Name: "Set Map and Pair With columns instances",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{loukoum.Column("foo"): 2, "a": 1}),
		String:     "UPDATE table SET a = 1, foo = 2",
		Query:      "UPDATE table SET a = $1, foo = $2",
		NamedQuery: "UPDATE table SET a = :arg_1, foo = :arg_2",
		Args:       []interface{}{1, 2},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
		},
	},
	{
		Name: "Set Map And Pair Uniqueness Variadic",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Pair("a", 1),
				loukoum.Map{"b": 2},
				loukoum.Pair("c", 3),
			),
		String:     "UPDATE table SET a = 1, b = 2, c = 3",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{1, 2, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
		},
	},
	{
		Name: "Set Map And Pair Uniqueness  Multiple Set() calls",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("a", 1)).
			Set(loukoum.Map{"b": 2}, loukoum.Pair("b", 2)).
			Set(loukoum.Map{"c": 3}, loukoum.Pair("c", 3)),
		String:     "UPDATE table SET a = 1, b = 2, c = 3",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{1, 2, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
		},
	},
	{
		Name: "Set Map And Pair Last Write Wins Variadic",
		Builder: loukoum.
			Update("table").
			Set(
				loukoum.Map{"a": 1, "b": 2},
				loukoum.Pair("a", 2),
				loukoum.Map{"b": 3},
				loukoum.Pair("c", 4),
			),
		String:     "UPDATE table SET a = 2, b = 3, c = 4",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{2, 3, 4},
		NamedArgs: map[string]interface{}{
			"arg_1": 2,
			"arg_2": 3,
			"arg_3": 4,
		},
	},
	{
		Name: "Set Map And Pair Last Write Wins Multiple Set() calls",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1, "b": 2}, loukoum.Pair("a", 2)).
			Set(loukoum.Map{"b": 2}, loukoum.Pair("b", 3)).
			Set(loukoum.Map{"c": 3}, loukoum.Pair("c", 4)),
		String:     "UPDATE table SET a = 2, b = 3, c = 4",
		Query:      "UPDATE table SET a = $1, b = $2, c = $3",
		NamedQuery: "UPDATE table SET a = :arg_1, b = :arg_2, c = :arg_3",
		Args:       []interface{}{2, 3, 4},
		NamedArgs: map[string]interface{}{
			"arg_1": 2,
			"arg_2": 3,
			"arg_3": 4,
		},
	},
	{
		Name:       "Set Valuer sql.NullTime valid",
		Builder:    loukoum.Update("table").Set(loukoum.Map{"created_at": pq.NullTime{Time: now, Valid: true}}),
		String:     fmt.Sprint("UPDATE table SET created_at = ", format.Time(now)),
		Query:      "UPDATE table SET created_at = $1",
		NamedQuery: "UPDATE table SET created_at = :arg_1",
		Args:       []interface{}{now},
		NamedArgs: map[string]interface{}{
			"arg_1": now,
		},
	},
	{
		Name:       "Set Valuer sql.NullTime null",
		Builder:    loukoum.Update("table").Set(loukoum.Map{"created_at": pq.NullTime{Time: now, Valid: false}}),
		String:     "UPDATE table SET created_at = NULL",
		Query:      "UPDATE table SET created_at = NULL",
		NamedQuery: "UPDATE table SET created_at = NULL",
	},
	{
		Name:       "Set Valuer sql.NullString valid",
		Builder:    loukoum.Update("table").Set(loukoum.Map{"message": sql.NullString{String: "ok", Valid: true}}),
		String:     "UPDATE table SET message = 'ok'",
		Query:      "UPDATE table SET message = $1",
		NamedQuery: "UPDATE table SET message = :arg_1",
		Args:       []interface{}{"ok"},
		NamedArgs: map[string]interface{}{
			"arg_1": "ok",
		},
	},
	{
		Name:       "Set Valuer sql.NullString null",
		Builder:    loukoum.Update("table").Set(loukoum.Map{"message": sql.NullString{String: "ok", Valid: false}}),
		String:     "UPDATE table SET message = NULL",
		Query:      "UPDATE table SET message = NULL",
		NamedQuery: "UPDATE table SET message = NULL",
	},
	{
		Name:       "Set Valuer sql.NullInt64 valid",
		Builder:    loukoum.Update("table").Set(loukoum.Map{"count": sql.NullInt64{Int64: 30, Valid: true}}),
		String:     "UPDATE table SET count = 30",
		Query:      "UPDATE table SET count = $1",
		NamedQuery: "UPDATE table SET count = :arg_1",
		Args:       []interface{}{int64(30)},
		NamedArgs: map[string]interface{}{
			"arg_1": int64(30),
		},
	},
	{
		Name:       "Set Valuer sql.NullInt64 null",
		Builder:    loukoum.Update("table").Set(loukoum.Map{"count": sql.NullInt64{Int64: 30, Valid: false}}),
		String:     "UPDATE table SET count = NULL",
		Query:      "UPDATE table SET count = NULL",
		NamedQuery: "UPDATE table SET count = NULL",
	},
	{
		Name: "Using without any arguments must panic",
		Failure: func() builder.Builder {
			return loukoum.Update("table").Set("a", "b", "c").Using()
		},
	},
	{
		Name: "Using() without any columns must panic",
		Failure: func() builder.Builder {
			return loukoum.Update("table").Set(loukoum.Pair("a", 30)).Using()
		},
	},
	{
		Name: "Set Using Muti-values A",
		Builder: loukoum.
			Update("table").
			Set("a", "b", "c").
			Using("d", "e", "f"),
		String:     "UPDATE table SET (a, b, c) = ('d', 'e', 'f')",
		Query:      "UPDATE table SET (a, b, c) = ($1, $2, $3)",
		NamedQuery: "UPDATE table SET (a, b, c) = (:arg_1, :arg_2, :arg_3)",
		Args:       []interface{}{"d", "e", "f"},
		NamedArgs: map[string]interface{}{
			"arg_1": "d",
			"arg_2": "e",
			"arg_3": "f",
		},
	},
	{
		Name: "Set Using Muti-values B",
		Builder: loukoum.
			Update("table").
			Set("a").
			Set("b").
			Set("c").
			Using("d", "e", "f"),
		String:     "UPDATE table SET (a, b, c) = ('d', 'e', 'f')",
		Query:      "UPDATE table SET (a, b, c) = ($1, $2, $3)",
		NamedQuery: "UPDATE table SET (a, b, c) = (:arg_1, :arg_2, :arg_3)",
		Args:       []interface{}{"d", "e", "f"},
		NamedArgs: map[string]interface{}{
			"arg_1": "d",
			"arg_2": "e",
			"arg_3": "f",
		},
	},
	{
		Name: "Set Using Columns to columns",
		Builder: loukoum.
			Update("table").
			Set("a", "b", "c").
			Using(loukoum.Raw("d+1"), loukoum.Raw("e+1"), loukoum.Raw("f+1")),
		String:     "UPDATE table SET (a, b, c) = (d+1, e+1, f+1)",
		Query:      "UPDATE table SET (a, b, c) = (d+1, e+1, f+1)",
		NamedQuery: "UPDATE table SET (a, b, c) = (d+1, e+1, f+1)",
	},
	{
		Name: "Set Using Sub-select",
		Builder: loukoum.Update("table").
			Set("a", "b", "c").
			Using(
				loukoum.Select("a", "b", "c").
					From("table").
					Where(loukoum.Condition("disabled").Equal(false)).Statement(),
			),
		String:     "UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = false))",
		Query:      "UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = $1))",
		NamedQuery: "UPDATE table SET (a, b, c) = (SELECT a, b, c FROM table WHERE (disabled = :arg_1))",
		Args:       []interface{}{false},
		NamedArgs: map[string]interface{}{
			"arg_1": false,
		},
	},
	{
		Name:       "Only Table",
		Builder:    loukoum.Update("table").Only().Set(loukoum.Map{"a": 1}),
		String:     "UPDATE ONLY table SET a = 1",
		Query:      "UPDATE ONLY table SET a = $1",
		NamedQuery: "UPDATE ONLY table SET a = :arg_1",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Simple Where",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)),
		String:     "UPDATE table SET a = 1 WHERE (id = 1)",
		Query:      "UPDATE table SET a = $1 WHERE (id = $2)",
		NamedQuery: "UPDATE table SET a = :arg_1 WHERE (id = :arg_2)",
		Args:       []interface{}{1, 1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 1,
		},
	},
	{
		Name: "Where And",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("status").Equal("online")),
		String:     "UPDATE table SET a = 1 WHERE ((id = 1) AND (status = 'online'))",
		Query:      "UPDATE table SET a = $1 WHERE ((id = $2) AND (status = $3))",
		NamedQuery: "UPDATE table SET a = :arg_1 WHERE ((id = :arg_2) AND (status = :arg_3))",
		Args:       []interface{}{1, 1, "online"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 1,
			"arg_3": "online",
		},
	},
	{
		Name: "Where And",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("status").Equal("online")),
		String:     "UPDATE table SET a = 1 WHERE ((id = 1) OR (status = 'online'))",
		Query:      "UPDATE table SET a = $1 WHERE ((id = $2) OR (status = $3))",
		NamedQuery: "UPDATE table SET a = :arg_1 WHERE ((id = :arg_2) OR (status = :arg_3))",
		Args:       []interface{}{1, 1, "online"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 1,
			"arg_3": "online",
		},
	},
	{
		Name: "From",
		Builder: loukoum.
			Update("table1").
			Set(loukoum.Map{"a": 1}).
			From("table2").
			Where(loukoum.Condition("table2.id").Equal(loukoum.Raw("table1.id"))),
		String:     "UPDATE table1 SET a = 1 FROM table2 WHERE (table2.id = table1.id)",
		Query:      "UPDATE table1 SET a = $1 FROM table2 WHERE (table2.id = table1.id)",
		NamedQuery: "UPDATE table1 SET a = :arg_1 FROM table2 WHERE (table2.id = table1.id)",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Returning",
		Builder: loukoum.
			Update("table").
			Set(loukoum.Map{"a": 1}).
			Returning("*"),
		String:     "UPDATE table SET a = 1 RETURNING *",
		Query:      "UPDATE table SET a = $1 RETURNING *",
		NamedQuery: "UPDATE table SET a = :arg_1 RETURNING *",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
}

func TestUpdate(t *testing.T) {
	for _, tt := range updatetests {
		t.Run(tt.Name, tt.Run)
	}
}
