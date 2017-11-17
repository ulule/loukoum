package loukoum_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestSelect(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Select("test")
		is.Equal("SELECT test", query.String())
	}
	{
		query := loukoum.SelectDistinct("test")
		is.Equal("SELECT DISTINCT test", query.String())
	}
	{
		query := loukoum.Select(loukoum.Column("test").As("foobar"))
		is.Equal("SELECT test AS foobar", query.String())
	}
	{
		query := loukoum.Select("test", "foobar")
		is.Equal("SELECT test, foobar", query.String())
	}
	{
		query := loukoum.Select("test", loukoum.Column("test2").As("foobar"))
		is.Equal("SELECT test, test2 AS foobar", query.String())
	}
	{
		query := loukoum.Select("a", "b", loukoum.Column("c").As("x"))
		is.Equal("SELECT a, b, c AS x", query.String())
	}
	{
		query := loukoum.Select("a", loukoum.Column("b"), loukoum.Column("c").As("x"))
		is.Equal("SELECT a, b, c AS x", query.String())
	}
}

func TestFrom(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Select("a", "b", "c").From("foobar")
		is.Equal("SELECT a, b, c FROM foobar", query.String())
	}
	{
		query := loukoum.Select("a").From(loukoum.Table("foobar").As("example"))
		is.Equal("SELECT a FROM foobar AS example", query.String())
	}
}

func TestJoin(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2 ON test1.id = test2.fk_id")

		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2", "test1.id = test2.fk_id")

		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test2", "test1.id = test2.fk_id", loukoum.InnerJoin)

		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test1").
			Join("test3", "test3.fkey = test1.id", loukoum.LeftJoin)

		is.Equal("SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", "test4.gid = test2.id", loukoum.RightJoin)

		is.Equal("SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test5").
			Join("test3", "ON test3.id = test5.fk_id", loukoum.InnerJoin)

		is.Equal("SELECT a, b, c FROM test5 INNER JOIN test3 ON test3.id = test5.fk_id", query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", "test4.gid = test2.id").Join("test3", "test4.uid = test3.id")

		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join("test4", loukoum.On("test4.gid", "test2.id")).
			Join("test3", loukoum.On("test4.uid", "test3.id"))

		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
	{
		query := loukoum.
			Select("a", "b", "c").
			From("test2").
			Join(loukoum.Table("test4"), loukoum.On("test4.gid", "test2.id")).
			Join(loukoum.Table("test3"), loukoum.On("test4.uid", "test3.id"))

		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
}

func TestWhere(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1))

		is.Equal(`SELECT id FROM table WHERE (id = 1)`, query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("'foo'"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (slug = 'foo'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("'foo'")).
			And(loukoum.Condition("title").Equal("'hello'"))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("slug").Equal("'foo'")).
			Or(loukoum.Condition("title").Equal("'hello'"))

		is.Equal("SELECT id FROM table WHERE (((id = 1) OR (slug = 'foo')) OR (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("slug").Equal("'foo'")).
			Or(loukoum.Condition("title").Equal("'hello'"))

		is.Equal(`SELECT id FROM table WHERE (((id = 1) OR (slug = 'foo')) OR (title = 'hello'))`, query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.Or(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("'hello'")),
			).
			Or(
				loukoum.And(loukoum.Condition("slug").Equal("'foo'"), loukoum.Condition("active").Equal(true)),
			)

		is.Equal(fmt.Sprint("SELECT id FROM table WHERE (((id = 1) OR (title = 'hello')) OR ",
			"((slug = 'foo') AND (active = true)))"), query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.And(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("'hello'")),
			)

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("'hello'"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("'hello'"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(
				loukoum.Condition("slug").Equal("'foo'").And(loukoum.Condition("active").Equal(true)),
			)

		is.Equal("SELECT id FROM table WHERE ((id = 1) OR ((slug = 'foo') AND (active = true)))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1).And(loukoum.Condition("slug").Equal("'foo'"))).
			Or(loukoum.Condition("active").Equal(true))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (active = true))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Equal(false))

		is.Equal("SELECT id FROM table WHERE (disabled = false)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").NotEqual(false))

		is.Equal("SELECT id FROM table WHERE (disabled != false)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Is(nil))

		is.Equal("SELECT id FROM table WHERE (disabled IS NULL)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("active").IsNot(true))

		is.Equal("SELECT id FROM table WHERE (active IS NOT true)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThan(2))

		is.Equal("SELECT id FROM table WHERE (count > 2)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThanOrEqual(4))

		is.Equal("SELECT id FROM table WHERE (count >= 4)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThan(3))

		is.Equal("SELECT id FROM table WHERE (count < 3)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThanOrEqual(6))

		is.Equal("SELECT id FROM table WHERE (count <= 6)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").Like("'foo%'"))

		is.Equal("SELECT id FROM table WHERE (title LIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotLike("'foo%'"))

		is.Equal("SELECT id FROM table WHERE (title NOT LIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").ILike("'foo%'"))

		is.Equal("SELECT id FROM table WHERE (title ILIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotILike("'foo%'"))

		is.Equal("SELECT id FROM table WHERE (title NOT ILIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").Between(10, 20))

		is.Equal("SELECT id FROM table WHERE (count BETWEEN 10 AND 20)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").NotBetween(10, 20))

		is.Equal("SELECT id FROM table WHERE (count NOT BETWEEN 10 AND 20)", query.String())
	}
}

func TestWhereIn(t *testing.T) {
	is := require.New(t)

	// Slice of integers
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In([]int64{1, 2, 3}))

		is.Equal("SELECT id FROM table WHERE (id IN (1, 2, 3))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn([]int{1, 2, 3}))

		is.Equal("SELECT id FROM table WHERE (id NOT IN (1, 2, 3))", query.String())
	}

	// Integers as variadic
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(1, 2, 3))

		is.Equal("SELECT id FROM table WHERE (id IN (1, 2, 3))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(1, 2, 3))

		is.Equal("SELECT id FROM table WHERE (id NOT IN (1, 2, 3))", query.String())
	}

	// Slice of strings
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"'read'", "'unread'"}))

		is.Equal("SELECT id FROM table WHERE (status IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"'read'", "'unread'"}))

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))", query.String())
	}

	// Strings as variadic
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("'read'", "'unread'"))

		is.Equal("SELECT id FROM table WHERE (status IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("'read'", "'unread'"))

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("'read'"))

		is.Equal("SELECT id FROM table WHERE (status IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("'read'"))

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read'))", query.String())
	}

	// Subquery
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Statement(),
			))

		is.Equal("SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = 1)))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)).
					Statement(),
			))

		is.Equal("SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = 1)))", query.String())
	}
}

// func TestInsert(t *testing.T) {
// 	is := require.New(t)

// 	{
// 		query := loukoum.
// 			Insert("table").
// 			Columns("a", "b", "c").
// 			Values([]string{"va", "vb", "vc"})

// 		is.Equal("INSERT INTO table (a, b, c) VALUES (va, vb, vc)", query.String())
// 	}
// }
