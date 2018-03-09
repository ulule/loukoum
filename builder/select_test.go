package builder_test

import (
	"fmt"
	"testing"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

var selecttests = []BuilderTest{
	{
		Name:       "Select",
		Builder:    loukoum.Select("test"),
		String:     "SELECT test",
		Query:      "SELECT test",
		NamedQuery: "SELECT test",
	},
	{
		Name:       "Distinct",
		Builder:    loukoum.Select("test").Distinct(),
		String:     "SELECT DISTINCT test",
		Query:      "SELECT DISTINCT test",
		NamedQuery: "SELECT DISTINCT test",
	},
	{
		Name:       "As",
		Builder:    loukoum.Select(loukoum.Column("test").As("foobar")),
		String:     "SELECT test AS foobar",
		Query:      "SELECT test AS foobar",
		NamedQuery: "SELECT test AS foobar",
	},
	{
		Name:       "Select two columns",
		Builder:    loukoum.Select("test", "foobar"),
		String:     "SELECT test, foobar",
		Query:      "SELECT test, foobar",
		NamedQuery: "SELECT test, foobar",
	},
	{
		Name:       "As with Column statement",
		Builder:    loukoum.Select("test", loukoum.Column("test2").As("foobar")),
		String:     "SELECT test, test2 AS foobar",
		Query:      "SELECT test, test2 AS foobar",
		NamedQuery: "SELECT test, test2 AS foobar",
	},
	{
		Name:       "As with mixed columns types A",
		Builder:    loukoum.Select("a", "b", loukoum.Column("c").As("x")),
		String:     "SELECT a, b, c AS x",
		Query:      "SELECT a, b, c AS x",
		NamedQuery: "SELECT a, b, c AS x",
	},
	{
		Name:       "As with mixed column types B",
		Builder:    loukoum.Select("a", loukoum.Column("b"), loukoum.Column("c").As("x")),
		String:     "SELECT a, b, c AS x",
		Query:      "SELECT a, b, c AS x",
		NamedQuery: "SELECT a, b, c AS x",
	},
	{
		Name:       "slice of column names",
		Builder:    loukoum.Select([]string{"a", "b", "c"}),
		String:     "SELECT a, b, c",
		Query:      "SELECT a, b, c",
		NamedQuery: "SELECT a, b, c",
	},
	{
		Name: "slice of column statements",
		Builder: loukoum.Select([]stmt.Column{
			loukoum.Column("a"),
			loukoum.Column("b"),
			loukoum.Column("c"),
		}),
		String:     "SELECT a, b, c",
		Query:      "SELECT a, b, c",
		NamedQuery: "SELECT a, b, c",
	},
	{
		Name:       "From",
		Builder:    loukoum.Select("a", "b", "c").From("foobar"),
		String:     "SELECT a, b, c FROM foobar",
		Query:      "SELECT a, b, c FROM foobar",
		NamedQuery: "SELECT a, b, c FROM foobar",
	},
	{
		Name:       "From column statement as",
		Builder:    loukoum.Select("a").From(loukoum.Table("foobar").As("example")),
		String:     "SELECT a FROM foobar AS example",
		Query:      "SELECT a FROM foobar AS example",
		NamedQuery: "SELECT a FROM foobar AS example",
	},
	{
		Name: "Join A",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2 ON test1.id = test2.fk_id"),
		String:     "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		Query:      "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		NamedQuery: "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
	},
	{
		Name: "Join B",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2", "test1.id = test2.fk_id"),
		String:     "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		Query:      "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		NamedQuery: "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
	},
	{
		Name: "Join inner",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2", "test1.id = test2.fk_id", loukoum.InnerJoin),
		String:     "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		Query:      "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
		NamedQuery: "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
	},
	{
		Name: "Join left",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test3", "test3.fkey = test1.id", loukoum.LeftJoin),
		String:     "SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id",
		Query:      "SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id",
		NamedQuery: "SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id",
	},
	{
		Name: "Join right",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", "test4.gid = test2.id", loukoum.RightJoin),
		String:     "SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id",
		Query:      "SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id",
		NamedQuery: "SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id",
	},
	{
		Name: "Join inner B",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test5").
			Join("test3", "ON test3.id = test5.fk_id", loukoum.InnerJoin),
		String:     "SELECT a, b, c FROM test5 INNER JOIN test3 ON test3.id = test5.fk_id",
		Query:      "SELECT a, b, c FROM test5 INNER JOIN test3 ON test3.id = test5.fk_id",
		NamedQuery: "SELECT a, b, c FROM test5 INNER JOIN test3 ON test3.id = test5.fk_id",
	},
	{
		Name: "Join Join",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", "test4.gid = test2.id").Join("test3", "test4.uid = test3.id"),
		String: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
		Query: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
		NamedQuery: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
	},
	{
		Name: "Join Join On",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", loukoum.On("test4.gid", "test2.id")).
			Join("test3", loukoum.On("test4.uid", "test3.id")),
		String: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
		Query: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
		NamedQuery: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
	},
	{
		Name: "Join Join with table statement",
		Builder: loukoum.
			Select("a", "b", "c").
			From("test2").
			Join(loukoum.Table("test4"), loukoum.On("test4.gid", "test2.id")).
			Join(loukoum.Table("test3"), loukoum.On("test4.uid", "test3.id")),
		String: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
		Query: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
		NamedQuery: fmt.Sprint(
			"SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id",
		),
	},
	{
		Name: "Where",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)),
		String:     `SELECT id FROM table WHERE (id = 1)`,
		Query:      `SELECT id FROM table WHERE (id = $1)`,
		NamedQuery: `SELECT id FROM table WHERE (id = :arg_1)`,
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Where two",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")),
		String:     "SELECT id FROM table WHERE ((id = 1) AND (slug = 'foo'))",
		Query:      "SELECT id FROM table WHERE ((id = $1) AND (slug = $2))",
		NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) AND (slug = :arg_2))",
		Args:       []interface{}{1, "foo"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
		},
	},
	{
		Name: "Where nested A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			And(loukoum.Condition("title").Equal("hello")),
		String:     "SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) AND (title = 'hello'))",
		Query:      "SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) AND (title = $3))",
		NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) AND (title = :arg_3))",
		Args:       []interface{}{1, "foo", "hello"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
			"arg_3": "hello",
		},
	},
	{
		Name: "Where nested B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			And(loukoum.Condition("title").Equal("hello")),
		String:     "SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) AND (title = 'hello'))",
		Query:      "SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) AND (title = $3))",
		NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) AND (title = :arg_3))",
		Args:       []interface{}{1, "foo", "hello"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
			"arg_3": "hello",
		},
	},
	{
		Name: "Where nested C",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("slug").Equal("foo")).
			Or(loukoum.Condition("title").Equal("hello")),
		String:     "SELECT id FROM table WHERE (((id = 1) OR (slug = 'foo')) OR (title = 'hello'))",
		Query:      "SELECT id FROM table WHERE (((id = $1) OR (slug = $2)) OR (title = $3))",
		NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) OR (slug = :arg_2)) OR (title = :arg_3))",
		Args:       []interface{}{1, "foo", "hello"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
			"arg_3": "hello",
		},
	},
	{
		Name: "Where nested D",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			Or(loukoum.Condition("title").Equal("hello")),
		String:     `SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (title = 'hello'))`,
		Query:      `SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) OR (title = $3))`,
		NamedQuery: `SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) OR (title = :arg_3))`,
		Args:       []interface{}{1, "foo", "hello"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
			"arg_3": "hello",
		},
	},
	{
		Name: "Where nested E",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.Or(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
			).
			Or(
				loukoum.And(loukoum.Condition("slug").Equal("foo"), loukoum.Condition("active").Equal(true)),
			),

		String: fmt.Sprint(
			"SELECT id FROM table WHERE (((id = 1) OR (title = 'hello')) OR ",
			"((slug = 'foo') AND (active = true)))",
		),
		Query: fmt.Sprint(
			"SELECT id FROM table WHERE (((id = $1) OR (title = $2)) OR ",
			"((slug = $3) AND (active = $4)))",
		),
		NamedQuery: fmt.Sprint(
			"SELECT id FROM table WHERE (((id = :arg_1) OR (title = :arg_2)) OR ",
			"((slug = :arg_3) AND (active = :arg_4)))",
		),
		Args: []interface{}{1, "hello", "foo", true},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "hello",
			"arg_3": "foo",
			"arg_4": true,
		},
	},
	{
		Name: "Where nested F",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.And(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
			),
		String:     "SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))",
		Query:      "SELECT id FROM table WHERE ((id = $1) AND (title = $2))",
		NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) AND (title = :arg_2))",
		Args:       []interface{}{1, "hello"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "hello",
		},
	},
	{
		Name: "Where nested G",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("hello")),
		String:     "SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))",
		Query:      "SELECT id FROM table WHERE ((id = $1) AND (title = $2))",
		NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) AND (title = :arg_2))",
		Args:       []interface{}{1, "hello"},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "hello",
		},
	},
	{
		Name: "Where nested H",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("hello")).
			Where(loukoum.Condition("disable").Equal(false)),
		String:     "SELECT id FROM table WHERE (((id = 1) AND (title = 'hello')) AND (disable = false))",
		Query:      "SELECT id FROM table WHERE (((id = $1) AND (title = $2)) AND (disable = $3))",
		NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (title = :arg_2)) AND (disable = :arg_3))",
		Args:       []interface{}{1, "hello", false},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "hello",
			"arg_3": false,
		},
	},
	{
		Name: "Where nested I",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(
				loukoum.Condition("slug").Equal("foo").And(loukoum.Condition("active").Equal(true)),
			),
		String:     "SELECT id FROM table WHERE ((id = 1) OR ((slug = 'foo') AND (active = true)))",
		Query:      "SELECT id FROM table WHERE ((id = $1) OR ((slug = $2) AND (active = $3)))",
		NamedQuery: "SELECT id FROM table WHERE ((id = :arg_1) OR ((slug = :arg_2) AND (active = :arg_3)))",
		Args:       []interface{}{1, "foo", true},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
			"arg_3": true,
		},
	},
	{
		Name: "Where nested J",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1).And(loukoum.Condition("slug").Equal("foo"))).
			Or(loukoum.Condition("active").Equal(true)),
		String:     "SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (active = true))",
		Query:      "SELECT id FROM table WHERE (((id = $1) AND (slug = $2)) OR (active = $3))",
		NamedQuery: "SELECT id FROM table WHERE (((id = :arg_1) AND (slug = :arg_2)) OR (active = :arg_3))",
		Args:       []interface{}{1, "foo", true},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": "foo",
			"arg_3": true,
		},
	},
	{
		Name: "Where Equal A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Equal(false)),
		String:     "SELECT id FROM table WHERE (disabled = false)",
		Query:      "SELECT id FROM table WHERE (disabled = $1)",
		NamedQuery: "SELECT id FROM table WHERE (disabled = :arg_1)",
		Args:       []interface{}{false},
		NamedArgs: map[string]interface{}{
			"arg_1": false,
		},
	},
	{
		Name: "Where Equal B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").NotEqual(false)),
		String:     "SELECT id FROM table WHERE (disabled != false)",
		Query:      "SELECT id FROM table WHERE (disabled != $1)",
		NamedQuery: "SELECT id FROM table WHERE (disabled != :arg_1)",
		Args:       []interface{}{false},
		NamedArgs: map[string]interface{}{
			"arg_1": false,
		},
	},
	{
		Name: "Where Is A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Is(nil)),
		String:     "SELECT id FROM table WHERE (disabled IS NULL)",
		Query:      "SELECT id FROM table WHERE (disabled IS NULL)",
		NamedQuery: "SELECT id FROM table WHERE (disabled IS NULL)",
	},
	{
		Name: "Where Is B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("active").IsNot(true)),
		String:     "SELECT id FROM table WHERE (active IS NOT true)",
		Query:      "SELECT id FROM table WHERE (active IS NOT $1)",
		NamedQuery: "SELECT id FROM table WHERE (active IS NOT :arg_1)",
		Args:       []interface{}{true},
		NamedArgs: map[string]interface{}{
			"arg_1": true,
		},
	},
	{
		Name: "Where Greater Than A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThan(2)),
		String:     "SELECT id FROM table WHERE (count > 2)",
		Query:      "SELECT id FROM table WHERE (count > $1)",
		NamedQuery: "SELECT id FROM table WHERE (count > :arg_1)",
		Args:       []interface{}{2},
		NamedArgs: map[string]interface{}{
			"arg_1": 2,
		},
	},
	{
		Name: "Where Greater Than B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThanOrEqual(4)),
		String:     "SELECT id FROM table WHERE (count >= 4)",
		Query:      "SELECT id FROM table WHERE (count >= $1)",
		NamedQuery: "SELECT id FROM table WHERE (count >= :arg_1)",
		Args:       []interface{}{4},
		NamedArgs: map[string]interface{}{
			"arg_1": 4,
		},
	},
	{
		Name: "Where Greater Than C",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("updated_at").GreaterThanOrEqual(loukoum.Raw("NOW()"))),
		String:     "SELECT id FROM table WHERE (updated_at >= NOW())",
		Query:      "SELECT id FROM table WHERE (updated_at >= NOW())",
		NamedQuery: "SELECT id FROM table WHERE (updated_at >= NOW())",
	},
	{
		Name: "Where Less Than A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThan(3)),
		String:     "SELECT id FROM table WHERE (count < 3)",
		Query:      "SELECT id FROM table WHERE (count < $1)",
		NamedQuery: "SELECT id FROM table WHERE (count < :arg_1)",
		Args:       []interface{}{3},
		NamedArgs: map[string]interface{}{
			"arg_1": 3,
		},
	},
	{
		Name: "Where Less Than B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThanOrEqual(6)),
		String:     "SELECT id FROM table WHERE (count <= 6)",
		Query:      "SELECT id FROM table WHERE (count <= $1)",
		NamedQuery: "SELECT id FROM table WHERE (count <= :arg_1)",
		Args:       []interface{}{6},
		NamedArgs: map[string]interface{}{
			"arg_1": 6,
		},
	},
	{
		Name: "Where Like A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").Like("foo%")),
		String:     "SELECT id FROM table WHERE (title LIKE 'foo%')",
		Query:      "SELECT id FROM table WHERE (title LIKE $1)",
		NamedQuery: "SELECT id FROM table WHERE (title LIKE :arg_1)",
		Args:       []interface{}{"foo%"},
		NamedArgs: map[string]interface{}{
			"arg_1": "foo%",
		},
	},
	{
		Name: "Where Like B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotLike("foo%")),
		String:     "SELECT id FROM table WHERE (title NOT LIKE 'foo%')",
		Query:      "SELECT id FROM table WHERE (title NOT LIKE $1)",
		NamedQuery: "SELECT id FROM table WHERE (title NOT LIKE :arg_1)",
		Args:       []interface{}{"foo%"},
		NamedArgs: map[string]interface{}{
			"arg_1": "foo%",
		},
	},
	{
		Name: "Where Like C",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").ILike("foo%")),
		String:     "SELECT id FROM table WHERE (title ILIKE 'foo%')",
		Query:      "SELECT id FROM table WHERE (title ILIKE $1)",
		NamedQuery: "SELECT id FROM table WHERE (title ILIKE :arg_1)",
		Args:       []interface{}{"foo%"},
		NamedArgs: map[string]interface{}{
			"arg_1": "foo%",
		},
	},
	{
		Name: "Where Like D",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotILike("foo%")),
		String:     "SELECT id FROM table WHERE (title NOT ILIKE 'foo%')",
		Query:      "SELECT id FROM table WHERE (title NOT ILIKE $1)",
		NamedQuery: "SELECT id FROM table WHERE (title NOT ILIKE :arg_1)",
		Args:       []interface{}{"foo%"},
		NamedArgs: map[string]interface{}{
			"arg_1": "foo%",
		},
	},
	{
		Name: "Where Between A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").Between(10, 20)),
		String:     "SELECT id FROM table WHERE (count BETWEEN 10 AND 20)",
		Query:      "SELECT id FROM table WHERE (count BETWEEN $1 AND $2)",
		NamedQuery: "SELECT id FROM table WHERE (count BETWEEN :arg_1 AND :arg_2)",
		Args:       []interface{}{10, 20},
		NamedArgs: map[string]interface{}{
			"arg_1": 10,
			"arg_2": 20,
		},
	},
	{
		Name: "Where Between B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").NotBetween(50, 70)),
		String:     "SELECT id FROM table WHERE (count NOT BETWEEN 50 AND 70)",
		Query:      "SELECT id FROM table WHERE (count NOT BETWEEN $1 AND $2)",
		NamedQuery: "SELECT id FROM table WHERE (count NOT BETWEEN :arg_1 AND :arg_2)",
		Args:       []interface{}{50, 70},
		NamedArgs: map[string]interface{}{
			"arg_1": 50,
			"arg_2": 70,
		},
	},
	{
		Name: "Where In Slice of integers",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In([]int64{1, 2, 3})),
		String:     "SELECT id FROM table WHERE (id IN (1, 2, 3))",
		Query:      "SELECT id FROM table WHERE (id IN ($1, $2, $3))",
		NamedQuery: "SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))",
		Args:       []interface{}{int64(1), int64(2), int64(3)},
		NamedArgs: map[string]interface{}{
			"arg_1": int64(1),
			"arg_2": int64(2),
			"arg_3": int64(3),
		},
	},
	{
		Name: "Where Not In Slice of integers",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn([]int{1, 2, 3})),
		String:     "SELECT id FROM table WHERE (id NOT IN (1, 2, 3))",
		Query:      "SELECT id FROM table WHERE (id NOT IN ($1, $2, $3))",
		NamedQuery: "SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))",
		Args:       []interface{}{1, 2, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
		},
	},
	{
		Name: "Where In integers variadic",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(1, 2, 3)),
		String:     "SELECT id FROM table WHERE (id IN (1, 2, 3))",
		Query:      "SELECT id FROM table WHERE (id IN ($1, $2, $3))",
		NamedQuery: "SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))",
		Args:       []interface{}{1, 2, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
		},
	},
	{
		Name: "Where Not In integers variadic",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(1, 2, 3)),
		String:     "SELECT id FROM table WHERE (id NOT IN (1, 2, 3))",
		Query:      "SELECT id FROM table WHERE (id NOT IN ($1, $2, $3))",
		NamedQuery: "SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))",
		Args:       []interface{}{1, 2, 3},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
			"arg_2": 2,
			"arg_3": 3,
		},
	},
	{
		Name: "Where In slice of strings",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"read", "unread"})),
		String:     "SELECT id FROM table WHERE (status IN ('read', 'unread'))",
		Query:      "SELECT id FROM table WHERE (status IN ($1, $2))",
		NamedQuery: "SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))",
		Args:       []interface{}{"read", "unread"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
			"arg_2": "unread",
		},
	},
	{
		Name: "Where Not In slice of strings",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"read", "unread"})),
		String:     "SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))",
		Query:      "SELECT id FROM table WHERE (status NOT IN ($1, $2))",
		NamedQuery: "SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))",
		Args:       []interface{}{"read", "unread"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
			"arg_2": "unread",
		},
	},
	{
		Name: "Where In strings as variadic",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("read", "unread")),
		String:     "SELECT id FROM table WHERE (status IN ('read', 'unread'))",
		Query:      "SELECT id FROM table WHERE (status IN ($1, $2))",
		NamedQuery: "SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))",
		Args:       []interface{}{"read", "unread"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
			"arg_2": "unread",
		},
	},
	{
		Name: "Where Not In strings as variadic",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("read", "unread")),
		String:     "SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))",
		Query:      "SELECT id FROM table WHERE (status NOT IN ($1, $2))",
		NamedQuery: "SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))",
		Args:       []interface{}{"read", "unread"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
			"arg_2": "unread",
		},
	},
	{
		Name: "Where In strings as variadic single",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("read")),
		String:     "SELECT id FROM table WHERE (status IN ('read'))",
		Query:      "SELECT id FROM table WHERE (status IN ($1))",
		NamedQuery: "SELECT id FROM table WHERE (status IN (:arg_1))",
		Args:       []interface{}{"read"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
		},
	},
	{
		Name: "Where Not In strings as variadic single",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("read")),
		String:     "SELECT id FROM table WHERE (status NOT IN ('read'))",
		Query:      "SELECT id FROM table WHERE (status NOT IN ($1))",
		NamedQuery: "SELECT id FROM table WHERE (status NOT IN (:arg_1))",
		Args:       []interface{}{"read"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
		},
	},
	{
		Name: "Where In strings as single slice",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"read"})),
		String:     "SELECT id FROM table WHERE (status IN ('read'))",
		Query:      "SELECT id FROM table WHERE (status IN ($1))",
		NamedQuery: "SELECT id FROM table WHERE (status IN (:arg_1))",
		Args:       []interface{}{"read"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
		},
	},
	{
		Name: "Where Not In strings as single slice",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"read"})),
		String:     "SELECT id FROM table WHERE (status NOT IN ('read'))",
		Query:      "SELECT id FROM table WHERE (status NOT IN ($1))",
		NamedQuery: "SELECT id FROM table WHERE (status NOT IN (:arg_1))",
		Args:       []interface{}{"read"},
		NamedArgs: map[string]interface{}{
			"arg_1": "read",
		},
	},
	{
		Name: "Where In subquery A",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Statement(),
			)),
		String:     "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = 1)))",
		Query:      "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = $1)))",
		NamedQuery: "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = :arg_1)))",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Where In subquery B",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Statement(),
			)),
		String:     "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = 1)))",
		Query:      "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = $1)))",
		NamedQuery: "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = :arg_1)))",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Where In subquery C",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)),
			)),
		String:     "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = 1)))",
		Query:      "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = $1)))",
		NamedQuery: "SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = :arg_1)))",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Where In subquery D",
		Builder: loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)),
			)),
		String:     "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = 1)))",
		Query:      "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = $1)))",
		NamedQuery: "SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = :arg_1)))",
		Args:       []interface{}{1},
		NamedArgs: map[string]interface{}{
			"arg_1": 1,
		},
	},
	{
		Name: "Where Group By one column",
		Builder: loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name"),
		String:     "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
		Query:      "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
		NamedQuery: "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
	},
	{
		Name: "Where Group By many columns A",
		Builder: loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy(loukoum.Column("name")),
		String:     "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
		Query:      "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
		NamedQuery: "SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name",
	},
	{
		Name: "Where Group By many columns B",
		Builder: loukoum.
			Select("name", "locale", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name", "locale"),
		String: fmt.Sprint(
			"SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
		),
		Query: fmt.Sprint(
			"SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
		),
		NamedQuery: fmt.Sprint(
			"SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
		),
	},
	{
		Name: "Where Group By many columns C",
		Builder: loukoum.
			Select("name", "locale", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy(loukoum.Column("name"), loukoum.Column("locale")),
		String: fmt.Sprint(
			"SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
		),
		Query: fmt.Sprint(
			"SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
		),
		NamedQuery: fmt.Sprint(
			"SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale",
		),
	},
	{
		Name: "Where Group By many columns D",
		Builder: loukoum.
			Select("name", "locale", "country", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name", "locale", "country"),
		String: fmt.Sprint(
			"SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
		),
		Query: fmt.Sprint(
			"SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
		),
		NamedQuery: fmt.Sprint(
			"SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
		),
	},
	{
		Name: "Where Group By many columns E",
		Builder: loukoum.
			Select("name", "locale", "country", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy(loukoum.Column("name"), loukoum.Column("locale"), loukoum.Column("country")),
		String: fmt.Sprint(
			"SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
		),
		Query: fmt.Sprint(
			"SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
		),
		NamedQuery: fmt.Sprint(
			"SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country",
		),
	},
	{
		Name: "Having one condition",
		Builder: loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name").
			Having(loukoum.Condition("COUNT(*)").GreaterThan(10)),
		String: fmt.Sprint(
			"SELECT name, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > 10)",
		),
		Query: fmt.Sprint(
			"SELECT name, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > $1)",
		),
		NamedQuery: fmt.Sprint(
			"SELECT name, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > :arg_1)",
		),
		Args: []interface{}{10},
		NamedArgs: map[string]interface{}{
			"arg_1": 10,
		},
	},
	{
		Name: "Having two conditions",
		Builder: loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name").
			Having(
				loukoum.Condition("COUNT(*)").GreaterThan(10).And(loukoum.Condition("COUNT(*)").LessThan(500)),
			),
		String: fmt.Sprint(
			"SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
			"HAVING ((COUNT(*) > 10) AND (COUNT(*) < 500))",
		),
		Query: fmt.Sprint(
			"SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
			"HAVING ((COUNT(*) > $1) AND (COUNT(*) < $2))",
		),
		NamedQuery: fmt.Sprint(
			"SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
			"HAVING ((COUNT(*) > :arg_1) AND (COUNT(*) < :arg_2))",
		),
		Args: []interface{}{10, 500},
		NamedArgs: map[string]interface{}{
			"arg_1": 10,
			"arg_2": 500,
		},
	},
	{
		Name: "Order By",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("id")),
		String:     "SELECT name FROM user ORDER BY id ASC",
		Query:      "SELECT name FROM user ORDER BY id ASC",
		NamedQuery: "SELECT name FROM user ORDER BY id ASC",
	},
	{
		Name: "Order By Asc",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("id", loukoum.Asc)),
		String:     "SELECT name FROM user ORDER BY id ASC",
		Query:      "SELECT name FROM user ORDER BY id ASC",
		NamedQuery: "SELECT name FROM user ORDER BY id ASC",
	},
	{
		Name: "Order By Desc",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("id", loukoum.Desc)),
		String:     "SELECT name FROM user ORDER BY id DESC",
		Query:      "SELECT name FROM user ORDER BY id DESC",
		NamedQuery: "SELECT name FROM user ORDER BY id DESC",
	},
	{
		Name: "Order By Multiple",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("locale"), loukoum.Order("id", loukoum.Desc)),
		String:     "SELECT name FROM user ORDER BY locale ASC, id DESC",
		Query:      "SELECT name FROM user ORDER BY locale ASC, id DESC",
		NamedQuery: "SELECT name FROM user ORDER BY locale ASC, id DESC",
	},
	{
		Name: "Order By Asc with column",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("id").Asc()),
		String:     "SELECT name FROM user ORDER BY id ASC",
		Query:      "SELECT name FROM user ORDER BY id ASC",
		NamedQuery: "SELECT name FROM user ORDER BY id ASC",
	},
	{
		Name: "Order By Desc with column",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("id").Desc()),
		String:     "SELECT name FROM user ORDER BY id DESC",
		Query:      "SELECT name FROM user ORDER BY id DESC",
		NamedQuery: "SELECT name FROM user ORDER BY id DESC",
	},
	{
		Name: "Order By Asc and Desc with column",
		Builder: loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("locale").Asc(), loukoum.Column("id").Desc()),
		String:     "SELECT name FROM user ORDER BY locale ASC, id DESC",
		Query:      "SELECT name FROM user ORDER BY locale ASC, id DESC",
		NamedQuery: "SELECT name FROM user ORDER BY locale ASC, id DESC",
	},
	{
		Name: "Limit with several types A",
		Builder: loukoum.
			Select("name").
			From("user").
			Limit(10),
		String:     "SELECT name FROM user LIMIT 10",
		Query:      "SELECT name FROM user LIMIT 10",
		NamedQuery: "SELECT name FROM user LIMIT 10",
	},
	{
		Name: "Limit with several types B",
		Builder: loukoum.
			Select("name").
			From("user").
			Limit("50"),
		String:     "SELECT name FROM user LIMIT 50",
		Query:      "SELECT name FROM user LIMIT 50",
		NamedQuery: "SELECT name FROM user LIMIT 50",
	},
	{
		Name: "Limit with several types C",
		Builder: loukoum.
			Select("name").
			From("user").
			Limit(uint64(700)),
		String:     "SELECT name FROM user LIMIT 700",
		Query:      "SELECT name FROM user LIMIT 700",
		NamedQuery: "SELECT name FROM user LIMIT 700",
	},
	{
		Name: "Limit corner case A",
		Failure: func() builder.Builder {
			return loukoum.Select("name").From("user").Limit(700.2)
		},
	},
	{
		Name: "Limit corner case B",
		Failure: func() builder.Builder {
			return loukoum.Select("name").From("user").Limit(float32(700.2))
		},
	},
	{
		Name: "Limit corner case C",
		Failure: func() builder.Builder {
			return loukoum.Select("name").From("user").Limit(-700)
		},
	},
	{
		Name: "Offset with several types A",
		Builder: loukoum.
			Select("name").
			From("user").
			Offset(10),
		String:     "SELECT name FROM user OFFSET 10",
		Query:      "SELECT name FROM user OFFSET 10",
		NamedQuery: "SELECT name FROM user OFFSET 10",
	},
	{
		Name: "Offset with several types B",
		Builder: loukoum.
			Select("name").
			From("user").
			Offset("50"),
		String:     "SELECT name FROM user OFFSET 50",
		Query:      "SELECT name FROM user OFFSET 50",
		NamedQuery: "SELECT name FROM user OFFSET 50",
	},
	{
		Name: "Offset with several types C",
		Builder: loukoum.
			Select("name").
			From("user").
			Offset(uint64(700)),
		String:     "SELECT name FROM user OFFSET 700",
		Query:      "SELECT name FROM user OFFSET 700",
		NamedQuery: "SELECT name FROM user OFFSET 700",
	},
	{
		Name: "Offset corner case A",
		Failure: func() builder.Builder {
			return loukoum.Select("name").From("user").Offset(700.2)
		},
	},
	{
		Name: "Offset corner case B",
		Failure: func() builder.Builder {
			return loukoum.Select("name").From("user").Offset(float32(700.2))
		},
	},
	{
		Name: "Offset corner case C",
		Failure: func() builder.Builder {
			return loukoum.Select("name").From("user").Offset(-700)
		},
	},
	{
		Name: "Prefix",
		Builder: loukoum.
			Select("name").
			From("user").
			Prefix("EXPLAIN ANALYZE"),
		String:     "EXPLAIN ANALYZE SELECT name FROM user",
		Query:      "EXPLAIN ANALYZE SELECT name FROM user",
		NamedQuery: "EXPLAIN ANALYZE SELECT name FROM user",
	},
	{
		Name: "Suffix",
		Builder: loukoum.
			Select("name").
			From("user").
			Suffix("FOR UPDATE"),
		String:     "SELECT name FROM user FOR UPDATE",
		Query:      "SELECT name FROM user FOR UPDATE",
		NamedQuery: "SELECT name FROM user FOR UPDATE",
	},
}

func TestSelect(t *testing.T) {
	for _, tt := range selecttests {
		t.Run(tt.Name, tt.Run)
	}
}
