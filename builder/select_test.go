package builder_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/stmt"
)

func TestSelect(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.Select("test")
		is.Equal("SELECT test", query.String())
	}
	{
		query := loukoum.Select("test").Distinct()
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
	{
		query := loukoum.Select([]string{"a", "b", "c"})
		is.Equal("SELECT a, b, c", query.String())
	}
	{
		query := loukoum.Select([]stmt.Column{loukoum.Column("a"), loukoum.Column("b"), loukoum.Column("c")})
		is.Equal("SELECT a, b, c", query.String())
	}
}

func TestSelect_From(t *testing.T) {
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

func TestSelect_Join(t *testing.T) {
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

func TestSelect_WhereOperatorOrder(t *testing.T) {
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
			And(loukoum.Condition("slug").Equal("foo"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (slug = 'foo'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			And(loukoum.Condition("title").Equal("hello"))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(loukoum.Condition("slug").Equal("foo")).
			Or(loukoum.Condition("title").Equal("hello"))

		is.Equal("SELECT id FROM table WHERE (((id = 1) OR (slug = 'foo')) OR (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			And(loukoum.Condition("slug").Equal("foo")).
			Or(loukoum.Condition("title").Equal("hello"))

		is.Equal(`SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (title = 'hello'))`, query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.Or(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
			).
			Or(
				loukoum.And(loukoum.Condition("slug").Equal("foo"), loukoum.Condition("active").Equal(true)),
			)

		is.Equal(fmt.Sprint("SELECT id FROM table WHERE (((id = 1) OR (title = 'hello')) OR ",
			"((slug = 'foo') AND (active = true)))"), query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(
				loukoum.And(loukoum.Condition("id").Equal(1), loukoum.Condition("title").Equal("hello")),
			)

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("hello"))

		is.Equal("SELECT id FROM table WHERE ((id = 1) AND (title = 'hello'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Where(loukoum.Condition("title").Equal("hello")).
			Where(loukoum.Condition("disable").Equal(false))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (title = 'hello')) AND (disable = false))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1)).
			Or(
				loukoum.Condition("slug").Equal("foo").And(loukoum.Condition("active").Equal(true)),
			)

		is.Equal("SELECT id FROM table WHERE ((id = 1) OR ((slug = 'foo') AND (active = true)))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").Equal(1).And(loukoum.Condition("slug").Equal("foo"))).
			Or(loukoum.Condition("active").Equal(true))

		is.Equal("SELECT id FROM table WHERE (((id = 1) AND (slug = 'foo')) OR (active = true))", query.String())
	}
}

func TestSelect_WhereEqual(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Equal(false))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (disabled = :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(false, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (disabled = false)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").NotEqual(false))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (disabled != :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(false, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (disabled != false)", query.String())
	}
}

func TestSelect_WhereIs(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("disabled").Is(nil))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (disabled IS NULL)", stmt)
		is.Len(args, 0)

		is.Equal("SELECT id FROM table WHERE (disabled IS NULL)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("active").IsNot(true))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (active IS NOT :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(true, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (active IS NOT true)", query.String())
	}
}

func TestSelect_WhereGreaterThan(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThan(2))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (count > :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(2, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (count > 2)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").GreaterThanOrEqual(4))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (count >= :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(4, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (count >= 4)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("updated_at").GreaterThanOrEqual(loukoum.Raw("NOW()")))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (updated_at >= NOW())", stmt)
		is.Len(args, 0)
		is.Equal("SELECT id FROM table WHERE (updated_at >= NOW())", query.String())
	}
}

func TestSelect_WhereLessThan(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThan(3))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (count < :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(3, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (count < 3)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").LessThanOrEqual(6))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (count <= :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal(6, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (count <= 6)", query.String())
	}
}

func TestSelect_WhereLike(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").Like("foo%"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (title LIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (title LIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotLike("foo%"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (title NOT LIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (title NOT LIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").ILike("foo%"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (title ILIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (title ILIKE 'foo%')", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("title").NotILike("foo%"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (title NOT ILIKE :arg_1)", stmt)
		is.Len(args, 1)
		is.Equal("foo%", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (title NOT ILIKE 'foo%')", query.String())
	}
}

func TestSelect_WhereBetween(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").Between(10, 20))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (count BETWEEN :arg_1 AND :arg_2)", stmt)
		is.Len(args, 2)
		is.Equal(10, args["arg_1"])
		is.Equal(20, args["arg_2"])

		is.Equal("SELECT id FROM table WHERE (count BETWEEN 10 AND 20)", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("count").NotBetween(50, 70))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (count NOT BETWEEN :arg_1 AND :arg_2)", stmt)
		is.Len(args, 2)
		is.Equal(50, args["arg_1"])
		is.Equal(70, args["arg_2"])

		is.Equal("SELECT id FROM table WHERE (count NOT BETWEEN 50 AND 70)", query.String())
	}
}

func TestSelect_WhereIn(t *testing.T) {
	is := require.New(t)

	// Slice of integers
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In([]int64{1, 2, 3}))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int64(1), args["arg_1"])
		is.Equal(int64(2), args["arg_2"])
		is.Equal(int64(3), args["arg_3"])

		is.Equal("SELECT id FROM table WHERE (id IN (1, 2, 3))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn([]int{1, 2, 3}))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int(1), args["arg_1"])
		is.Equal(int(2), args["arg_2"])
		is.Equal(int(3), args["arg_3"])

		is.Equal("SELECT id FROM table WHERE (id NOT IN (1, 2, 3))", query.String())
	}

	// Integers as variadic
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").In(1, 2, 3))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (id IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int(1), args["arg_1"])
		is.Equal(int(2), args["arg_2"])
		is.Equal(int(3), args["arg_3"])

		is.Equal("SELECT id FROM table WHERE (id IN (1, 2, 3))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(1, 2, 3))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (id NOT IN (:arg_1, :arg_2, :arg_3))", stmt)
		is.Len(args, 3)
		is.Equal(int(1), args["arg_1"])
		is.Equal(int(2), args["arg_2"])
		is.Equal(int(3), args["arg_3"])

		is.Equal("SELECT id FROM table WHERE (id NOT IN (1, 2, 3))", query.String())
	}

	// Slice of strings
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"read", "unread"}))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args["arg_1"])
		is.Equal("unread", args["arg_2"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"read", "unread"}))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args["arg_1"])
		is.Equal("unread", args["arg_2"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))", query.String())
	}

	// Strings as variadic
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("read", "unread"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args["arg_1"])
		is.Equal("unread", args["arg_2"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("read", "unread"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1, :arg_2))", stmt)
		is.Len(args, 2)
		is.Equal("read", args["arg_1"])
		is.Equal("unread", args["arg_2"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read', 'unread'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In("read"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn("read"))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (status NOT IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").In([]string{"read"}))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (status IN ('read'))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("status").NotIn([]string{"read"}))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (status NOT IN (:arg_1))", stmt)
		is.Len(args, 1)
		is.Equal("read", args["arg_1"])

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
					Where(loukoum.Condition("id").Equal(1)),
			))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = :arg_1)))", stmt)
		is.Len(args, 1)
		is.Equal(1, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (id IN (SELECT id FROM table WHERE (id = 1)))", query.String())
	}
	{
		query := loukoum.
			Select("id").
			From("table").
			Where(loukoum.Condition("id").NotIn(
				loukoum.Select("id").
					From("table").
					Where(loukoum.Condition("id").Equal(1)),
			))

		stmt, args := query.NamedQuery()
		is.Equal("SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = :arg_1)))", stmt)
		is.Len(args, 1)
		is.Equal(1, args["arg_1"])

		is.Equal("SELECT id FROM table WHERE (id NOT IN (SELECT id FROM table WHERE (id = 1)))", query.String())
	}
}

func TestSelect_GroupBy(t *testing.T) {
	is := require.New(t)

	// One column
	{
		query := loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name")

		is.Equal("SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name", query.String())
	}
	{
		query := loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy(loukoum.Column("name"))

		is.Equal("SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name", query.String())
	}

	// Many columns
	{
		query := loukoum.
			Select("name", "locale", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name", "locale")

		is.Equal(fmt.Sprint("SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale"), query.String())
	}
	{
		query := loukoum.
			Select("name", "locale", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy(loukoum.Column("name"), loukoum.Column("locale"))

		is.Equal(fmt.Sprint("SELECT name, locale, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale"), query.String())
	}
	{
		query := loukoum.
			Select("name", "locale", "country", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name", "locale", "country")

		is.Equal(fmt.Sprint("SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country"), query.String())
	}
	{
		query := loukoum.
			Select("name", "locale", "country", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy(loukoum.Column("name"), loukoum.Column("locale"), loukoum.Column("country"))

		is.Equal(fmt.Sprint("SELECT name, locale, country, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name, locale, country"), query.String())
	}
}

func TestSelect_Having(t *testing.T) {
	is := require.New(t)

	// One condition
	{
		query := loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name").
			Having(loukoum.Condition("COUNT(*)").GreaterThan(10))

		is.Equal(fmt.Sprint("SELECT name, COUNT(*) FROM user ",
			"WHERE (disabled IS NOT NULL) GROUP BY name HAVING (COUNT(*) > 10)"), query.String())
	}

	// Two conditions
	{
		query := loukoum.
			Select("name", "COUNT(*)").
			From("user").
			Where(loukoum.Condition("disabled").IsNull(false)).
			GroupBy("name").
			Having(
				loukoum.Condition("COUNT(*)").GreaterThan(10).And(loukoum.Condition("COUNT(*)").LessThan(500)),
			)

		is.Equal(fmt.Sprint("SELECT name, COUNT(*) FROM user WHERE (disabled IS NOT NULL) GROUP BY name ",
			"HAVING ((COUNT(*) > 10) AND (COUNT(*) < 500))"), query.String())
	}
}

func TestSelect_OrderBy(t *testing.T) {
	is := require.New(t)

	// With Order
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("id"))

		is.Equal("SELECT name FROM user ORDER BY id ASC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("id", loukoum.Asc))

		is.Equal("SELECT name FROM user ORDER BY id ASC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("id", loukoum.Desc))

		is.Equal("SELECT name FROM user ORDER BY id DESC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("locale"), loukoum.Order("id", loukoum.Desc))

		is.Equal("SELECT name FROM user ORDER BY locale ASC, id DESC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Order("locale")).
			OrderBy(loukoum.Order("id", loukoum.Desc))

		is.Equal("SELECT name FROM user ORDER BY locale ASC, id DESC", query.String())
	}

	// With Column
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("id").Asc())

		is.Equal("SELECT name FROM user ORDER BY id ASC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("id").Desc())

		is.Equal("SELECT name FROM user ORDER BY id DESC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("locale").Asc(), loukoum.Column("id").Desc())

		is.Equal("SELECT name FROM user ORDER BY locale ASC, id DESC", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			OrderBy(loukoum.Column("locale").Asc()).
			OrderBy(loukoum.Column("id").Desc())

		is.Equal("SELECT name FROM user ORDER BY locale ASC, id DESC", query.String())
	}
}

func TestSelect_Limit(t *testing.T) {
	is := require.New(t)

	// Limit with several types
	{
		query := loukoum.
			Select("name").
			From("user").
			Limit(10)

		is.Equal("SELECT name FROM user LIMIT 10", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			Limit("50")

		is.Equal("SELECT name FROM user LIMIT 50", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			Limit(uint64(700))

		is.Equal("SELECT name FROM user LIMIT 700", query.String())
	}

	// Corner cases...
	{
		Failure(is, func() builder.Builder {
			return loukoum.Select("name").From("user").Limit(700.2)
		})
		Failure(is, func() builder.Builder {
			return loukoum.Select("name").From("user").Limit(float32(700.2))
		})
		Failure(is, func() builder.Builder {
			return loukoum.Select("name").From("user").Limit(-700)
		})
	}
}

func TestSelect_Offset(t *testing.T) {
	is := require.New(t)

	// Offset with several types
	{
		query := loukoum.
			Select("name").
			From("user").
			Offset(10)

		is.Equal("SELECT name FROM user OFFSET 10", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			Offset("50")

		is.Equal("SELECT name FROM user OFFSET 50", query.String())
	}
	{
		query := loukoum.
			Select("name").
			From("user").
			Offset(uint64(700))

		is.Equal("SELECT name FROM user OFFSET 700", query.String())
	}

	// Corner cases...
	{
		Failure(is, func() builder.Builder {
			return loukoum.Select("name").From("user").Offset(700.2)
		})
		Failure(is, func() builder.Builder {
			return loukoum.Select("name").From("user").Offset(float32(700.2))
		})
		Failure(is, func() builder.Builder {
			return loukoum.Select("name").From("user").Offset(-700)
		})
	}
}

func TestSelect_Prefix(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("name").
			From("user").
			Prefix("EXPLAIN ANALYZE")

		is.Equal("EXPLAIN ANALYZE SELECT name FROM user", query.String())
	}
}

func TestSelect_Suffix(t *testing.T) {
	is := require.New(t)

	{
		query := loukoum.
			Select("name").
			From("user").
			Suffix("FOR UPDATE")

		is.Equal("SELECT name FROM user FOR UPDATE", query.String())
	}
}
