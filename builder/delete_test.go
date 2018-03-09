package builder_test

import (
	"fmt"
	"testing"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/format"
)

var deletetests = []BuilderTest{
	{
		Name:       "Simple",
		Builder:    loukoum.Delete("table"),
		String:     "DELETE FROM table",
		Query:      "DELETE FROM table",
		NamedQuery: "DELETE FROM table",
	},
	{
		Name:       "Only",
		Builder:    loukoum.Delete("table").Only(),
		String:     "DELETE FROM ONLY table",
		Query:      "DELETE FROM ONLY table",
		NamedQuery: "DELETE FROM ONLY table",
	},
	{
		Name:       "Table statement",
		Builder:    loukoum.Delete(loukoum.Table("table")),
		String:     "DELETE FROM table",
		Query:      "DELETE FROM table",
		NamedQuery: "DELETE FROM table",
	},
	{
		Name:       "As",
		Builder:    loukoum.Delete(loukoum.Table("table").As("foobar")),
		String:     "DELETE FROM table AS foobar",
		Query:      "DELETE FROM table AS foobar",
		NamedQuery: "DELETE FROM table AS foobar",
	},
	{
		Name:       "As and Only",
		Builder:    loukoum.Delete(loukoum.Table("table").As("foobar")).Only(),
		String:     "DELETE FROM ONLY table AS foobar",
		Query:      "DELETE FROM ONLY table AS foobar",
		NamedQuery: "DELETE FROM ONLY table AS foobar",
	},
	{
		Name:       "Using one table",
		Builder:    loukoum.Delete("table").Using("foobar"),
		String:     "DELETE FROM table USING foobar",
		Query:      "DELETE FROM table USING foobar",
		NamedQuery: "DELETE FROM table USING foobar",
	},
	{
		Name:       "Using one table statement",
		Builder:    loukoum.Delete("table").Using(loukoum.Table("foobar")),
		String:     "DELETE FROM table USING foobar",
		Query:      "DELETE FROM table USING foobar",
		NamedQuery: "DELETE FROM table USING foobar",
	},
	{
		Name:       "Using one table statement as",
		Builder:    loukoum.Delete("table").Using(loukoum.Table("foobar").As("foo")),
		String:     "DELETE FROM table USING foobar AS foo",
		Query:      "DELETE FROM table USING foobar AS foo",
		NamedQuery: "DELETE FROM table USING foobar AS foo",
	},
	{
		Name:       "Using two tables",
		Builder:    loukoum.Delete("table").Using("foobar", "example"),
		String:     "DELETE FROM table USING foobar, example",
		Query:      "DELETE FROM table USING foobar, example",
		NamedQuery: "DELETE FROM table USING foobar, example",
	},
	{
		Name:       "Using two tables with table statement",
		Builder:    loukoum.Delete("table").Using(loukoum.Table("foobar"), "example"),
		String:     "DELETE FROM table USING foobar, example",
		Query:      "DELETE FROM table USING foobar, example",
		NamedQuery: "DELETE FROM table USING foobar, example",
	},
	{
		Name:       "Using two tables with as",
		Builder:    loukoum.Delete("table").Using(loukoum.Table("example"), loukoum.Table("foobar").As("foo")),
		String:     "DELETE FROM table USING example, foobar AS foo",
		Query:      "DELETE FROM table USING example, foobar AS foo",
		NamedQuery: "DELETE FROM table USING example, foobar AS foo",
	},
	{
		Name:       "Where",
		Builder:    loukoum.Delete("table").Where(loukoum.Condition("id").Equal(1)),
		String:     "DELETE FROM table WHERE (id = 1)",
		Query:      "DELETE FROM table WHERE (id = $1)",
		NamedQuery: "DELETE FROM table WHERE (id = :arg_1)",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Where time",
		Builder: loukoum.Delete("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("created_at").GreaterThan(now)),
		String:     fmt.Sprintf("DELETE FROM table WHERE ((id = 1) AND (created_at > %s))", format.Time(now)),
		Query:      "DELETE FROM table WHERE ((id = $1) AND (created_at > $2))",
		NamedQuery: "DELETE FROM table WHERE ((id = :arg_1) AND (created_at > :arg_2))",
		Args:       []interface{}{1, now},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": now,
		},
	},
	{
		Name:       "Returning *",
		Builder:    loukoum.Delete("table").Returning("*"),
		String:     "DELETE FROM table RETURNING *",
		Query:      "DELETE FROM table RETURNING *",
		NamedQuery: "DELETE FROM table RETURNING *",
	},
	{
		Name:       "Returning id",
		Builder:    loukoum.Delete("table").Returning("id"),
		String:     "DELETE FROM table RETURNING id",
		Query:      "DELETE FROM table RETURNING id",
		NamedQuery: "DELETE FROM table RETURNING id",
	},
}

func TestDelete(t *testing.T) {
	for _, tt := range deletetests {
		t.Run(tt.Name, tt.Run)
	}
}
