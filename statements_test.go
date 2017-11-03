package loukoum_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestStatementSelect(t *testing.T) {
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
		query := loukoum.Select("a", "b", "c").From("foobar")
		is.Equal("SELECT a, b, c FROM foobar", query.String())
	}
	{
		query := loukoum.Select("a", "b", loukoum.Column("c").As("x")).From("foobar")
		is.Equal("SELECT a, b, c AS x FROM foobar", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test1").Join("test2 ON test1.id = test2.fk_id")
		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test1").Join("test2", "test1.id = test2.fk_id")
		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test1").
			Join("test2", "test1.id = test2.fk_id", loukoum.InnerJoin)
		is.Equal("SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test1").
			Join("test3", "test3.fkey = test1.id", loukoum.LeftJoin)
		is.Equal("SELECT a, b, c FROM test1 LEFT JOIN test3 ON test3.fkey = test1.id", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test2").
			Join("test4", "test4.gid = test2.id", loukoum.RightJoin)
		is.Equal("SELECT a, b, c FROM test2 RIGHT JOIN test4 ON test4.gid = test2.id", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test5").
			Join("test3", "ON test3.id = test5.fk_id", loukoum.InnerJoin)
		is.Equal("SELECT a, b, c FROM test5 INNER JOIN test3 ON test3.id = test5.fk_id", query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test2").
			Join("test4", "test4.gid = test2.id").Join("test3", "test4.uid = test3.id")
		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
	{
		query := loukoum.Select("a", "b", "c").From("test2").
			Join("test4", loukoum.On("test4.gid", "test2.id")).
			Join("test3", loukoum.On("test4.uid", "test3.id"))
		is.Equal(fmt.Sprint("SELECT a, b, c FROM test2 INNER JOIN test4 ON test4.gid = test2.id ",
			"INNER JOIN test3 ON test4.uid = test3.id"), query.String())
	}
}
