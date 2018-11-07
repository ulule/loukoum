package builder_test

import (
	"testing"
	"time"

	"github.com/ulule/loukoum/v2"
	"github.com/ulule/loukoum/v2/builder"
)

func TestDelete(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Simple",
			Builders: []builder.Builder{
				loukoum.Delete("table"),
				loukoum.Delete(loukoum.Table("table")),
			},
			SameQuery: "DELETE FROM table",
		},
		{
			Name: "Only",
			Builders: []builder.Builder{
				loukoum.Delete("table").Only(),
				loukoum.Delete(loukoum.Table("table")).Only(),
			},
			SameQuery: "DELETE FROM ONLY table",
		},
		{
			Name:      "As",
			Builder:   loukoum.Delete(loukoum.Table("table").As("foobar")),
			SameQuery: "DELETE FROM table AS foobar",
		},
		{
			Name:      "As only",
			Builder:   loukoum.Delete(loukoum.Table("table").As("foobar")).Only(),
			SameQuery: "DELETE FROM ONLY table AS foobar",
		},
	})
}

func TestDelete_Using(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "One table",
			Builders: []builder.Builder{
				loukoum.Delete("table").Using("foobar"),
				loukoum.Delete("table").Using(loukoum.Table("foobar")),
			},
			SameQuery: "DELETE FROM table USING foobar",
		},
		{
			Name:      "One table as",
			Builder:   loukoum.Delete("table").Using(loukoum.Table("foobar").As("foo")),
			SameQuery: "DELETE FROM table USING foobar AS foo",
		},
		{
			Name: "Two tables",
			Builders: []builder.Builder{
				loukoum.Delete("table").Using("foobar", "example"),
				loukoum.Delete("table").Using(loukoum.Table("foobar"), "example"),
			},
			SameQuery: "DELETE FROM table USING foobar, example",
		},
		{
			Name:      "Two tables as",
			Builder:   loukoum.Delete("table").Using(loukoum.Table("example"), loukoum.Table("foobar").As("foo")),
			SameQuery: "DELETE FROM table USING example, foobar AS foo",
		},
	})
}

func TestDelete_Where(t *testing.T) {
	when, err := time.Parse(time.RFC3339, "2017-11-23T17:47:27+01:00")
	if err != nil {
		t.Fatal(err)
	}
	RunBuilderTests(t, []BuilderTest{
		{
			Name:       "Simple",
			Builder:    loukoum.Delete("table").Where(loukoum.Condition("id").Equal(1)),
			String:     "DELETE FROM table WHERE (id = 1)",
			Query:      "DELETE FROM table WHERE (id = $1)",
			NamedQuery: "DELETE FROM table WHERE (id = :arg_1)",
			Args:       []interface{}{1},
		},
		{
			Name: "Complex",
			Builder: loukoum.Delete("table").
				Where(loukoum.Condition("id").Equal(1)).
				And(loukoum.Condition("created_at").GreaterThan(when)),
			String:     "DELETE FROM table WHERE ((id = 1) AND (created_at > '2017-11-23 16:47:27+00'))",
			Query:      "DELETE FROM table WHERE ((id = $1) AND (created_at > $2))",
			NamedQuery: "DELETE FROM table WHERE ((id = :arg_1) AND (created_at > :arg_2))",
			Args:       []interface{}{1, when},
		},
	})
}

func TestDelete_Returning(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:      "*",
			Builder:   loukoum.Delete("table").Returning("*"),
			SameQuery: "DELETE FROM table RETURNING *",
		},
		{
			Name:      "id",
			Builder:   loukoum.Delete("table").Returning("id"),
			SameQuery: "DELETE FROM table RETURNING id",
		},
	})
}
