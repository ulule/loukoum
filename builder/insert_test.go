package builder_test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"

	loukoum "github.com/ulule/loukoum/v3"
	"github.com/ulule/loukoum/v3/builder"
)

func TestInsert_Columns(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:      "With columns",
			Builder:   loukoum.Insert("table").Columns("a", "b", "c"),
			SameQuery: "INSERT INTO table (a, b, c)",
		},
		{
			Name:      "Without columns",
			Builder:   loukoum.Insert("table"),
			SameQuery: "INSERT INTO table",
		},
	})
}

func TestInsert_Comment(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name:      "With columns",
			Builder:   loukoum.Insert("table").Columns("a", "b", "c").Comment("/foo"),
			SameQuery: "INSERT INTO table (a, b, c); -- /foo",
		},
		{
			Name:      "Without columns",
			Builder:   loukoum.Insert("table").Comment("/foo"),
			SameQuery: "INSERT INTO table; -- /foo",
		},
	})
}

func TestInsert_Values(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "With columns",
			Builders: []builder.Builder{
				loukoum.Insert("table").Columns("a", "b", "c").Values([]string{"va", "vb", "vc"}),
				loukoum.Insert("table").Columns("a", "b", "c").Values("va", "vb", "vc"),
			},
			String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc')",
			Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3)",
			NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3)",
			Args:       []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "Without columns",
			Builders: []builder.Builder{
				loukoum.Insert("table").Values([]string{"va", "vb", "vc"}),
				loukoum.Insert("table").Values("va", "vb", "vc"),
			},
			String:     "INSERT INTO table VALUES ('va', 'vb', 'vc')",
			Query:      "INSERT INTO table VALUES ($1, $2, $3)",
			NamedQuery: "INSERT INTO table VALUES (:arg_1, :arg_2, :arg_3)",
			Args:       []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "With raw values",
			Builders: []builder.Builder{
				loukoum.Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")),
				loukoum.Insert("table").
					Columns("email", "enabled", "created_at").
					Values([]interface{}{"tech@ulule.com", true, loukoum.Raw("NOW()")}),
			},
			String:     "INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW())",
			Query:      "INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW())",
			NamedQuery: "INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW())",
			Args:       []interface{}{"tech@ulule.com", true},
		},
		{
			Name: "With byte slice",
			Builders: []builder.Builder{
				loukoum.Insert("table").
					Set(loukoum.Pair("data", []byte{1, 2, 3})),
				loukoum.Insert("table").
					Columns("data").
					Values([]byte{1, 2, 3}),
				loukoum.Insert("table").
					Columns("data").
					Values([][]byte{{1, 2, 3}}),
			},
			String:     "INSERT INTO table (data) VALUES (decode('010203', 'hex'))",
			Query:      "INSERT INTO table (data) VALUES ($1)",
			NamedQuery: "INSERT INTO table (data) VALUES (:arg_1)",
			Args:       []interface{}{[]byte{1, 2, 3}},
		},
	})
}

func TestInsert_OnConflict(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Do nothing without target",
			Builder: loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict(loukoum.DoNothing()),
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT DO NOTHING",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
				"ON CONFLICT DO NOTHING",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
				"ON CONFLICT DO NOTHING",
			),
			Args: []interface{}{"tech@ulule.com", true},
		},
		{
			Name: "Do nothing",
			Builders: []builder.Builder{
				loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict("email", loukoum.DoNothing()),
				loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict(loukoum.Column("email"), loukoum.DoNothing()),
			},
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT (email) DO NOTHING",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
				"ON CONFLICT (email) DO NOTHING",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
				"ON CONFLICT (email) DO NOTHING",
			),
			Args: []interface{}{"tech@ulule.com", true},
		},
		{
			Name: "Do nothing with multiple targets",
			Builder: loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", loukoum.Column("uuid"), "reference", loukoum.DoNothing()),
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT (email, uuid, reference) DO NOTHING",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
				"ON CONFLICT (email, uuid, reference) DO NOTHING",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
				"ON CONFLICT (email, uuid, reference) DO NOTHING",
			),
			Args: []interface{}{"tech@ulule.com", true},
		},
		{
			Name: "Do update",
			Builders: []builder.Builder{
				loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict("email", loukoum.DoUpdate(
						loukoum.Pair("created_at", loukoum.Raw("NOW()")),
						loukoum.Pair("enabled", true),
					)),
				loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict(loukoum.Column("email"), loukoum.DoUpdate(
						loukoum.Pair("created_at", loukoum.Raw("NOW()")),
						loukoum.Pair("enabled", true),
					)),
			},
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT (email) DO UPDATE SET created_at = NOW(), enabled = true",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
				"ON CONFLICT (email) DO UPDATE SET created_at = NOW(), enabled = $3",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
				"ON CONFLICT (email) DO UPDATE SET created_at = NOW(), enabled = :arg_3",
			),
			Args: []interface{}{"tech@ulule.com", true, true},
		},
		{
			Name: "Do update with two targets",
			Builder: loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", "uuid", loukoum.DoUpdate(
					loukoum.Pair("created_at", loukoum.Raw("NOW()")),
					loukoum.Pair("enabled", true),
				)),
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT (email, uuid) DO UPDATE SET created_at = NOW(), enabled = true",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
				"ON CONFLICT (email, uuid) DO UPDATE SET created_at = NOW(), enabled = $3",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
				"ON CONFLICT (email, uuid) DO UPDATE SET created_at = NOW(), enabled = :arg_3",
			),
			Args: []interface{}{"tech@ulule.com", true, true},
		},
		{
			Name: "Do update with 3 targets",
			Builder: loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", loukoum.Column("uuid"), "reference", loukoum.DoUpdate(
					loukoum.Pair("created_at", loukoum.Raw("NOW()")),
					loukoum.Pair("enabled", true),
				)),
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT (email, uuid, reference) DO UPDATE SET created_at = NOW(), enabled = true",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
				"ON CONFLICT (email, uuid, reference) DO UPDATE SET created_at = NOW(), enabled = $3",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
				"ON CONFLICT (email, uuid, reference) DO UPDATE SET created_at = NOW(), enabled = :arg_3",
			),
			Args: []interface{}{"tech@ulule.com", true, true},
		},
		{
			Name: "Corner case 0",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict()
			},
		},
		{
			Name: "Corner case 1",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict("email")
			},
		},
		{
			Name: "Corner case 2",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict(loukoum.DoUpdate(
						loukoum.Pair("created_at", loukoum.Raw("NOW()")),
						loukoum.Pair("enabled", true),
					))
			},
		},
		{
			Name: "Corner case 3",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict("email", 6700)
			},
		},
		{
			Name: "Corner case 4",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict(569)
			},
		},
		{
			Name: "Corner case 5",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict("email", "uuid")
			},
		},
		{
			Name: "Corner case 6",
			Failure: func() builder.Builder {
				return loukoum.
					Insert("table").
					Columns("email", "enabled", "created_at").
					Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
					OnConflict(loukoum.Column("email"), loukoum.Column("uuid"), loukoum.Column("reference"))
			},
		},
	})
}

func TestInsert_Returning(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "One column",
			Builder: loukoum.
				Insert("table").
				Columns("a", "b", "c").
				Values([]string{"va", "vb", "vc"}).
				Returning("a"),
			String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a",
			Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a",
			NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a",
			Args:       []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "Two columns",
			Builder: loukoum.
				Insert("table").
				Columns("a", "b", "c").
				Values([]string{"va", "vb", "vc"}).
				Returning("a", "b"),
			String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a, b",
			Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a, b",
			NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a, b",
			Args:       []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "Three columns",
			Builder: loukoum.
				Insert("table").
				Columns("a", "b", "c").
				Values([]string{"va", "vb", "vc"}).
				Returning("a", "b", "c"),
			String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a, b, c",
			Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a, b, c",
			NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a, b, c",
			Args:       []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "With alias",
			Builder: loukoum.
				Insert("table").
				Columns("a", "b", "c").
				Values([]string{"va", "vb", "vc"}).
				Returning(loukoum.Column("a").As("alias_a")),
			String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a AS alias_a",
			Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a AS alias_a",
			NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a AS alias_a",
			Args:       []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "With two aliases",
			Builder: loukoum.
				Insert("table").
				Columns("a", "b", "c").
				Values([]string{"va", "vb", "vc"}).
				Returning(loukoum.Column("a").As("alias_a"), loukoum.Column("b").As("alias_b")),
			String: fmt.Sprint(
				"INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') ",
				"RETURNING a AS alias_a, b AS alias_b",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (a, b, c) VALUES ($1, $2, $3) ",
				"RETURNING a AS alias_a, b AS alias_b",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) ",
				"RETURNING a AS alias_a, b AS alias_b",
			),
			Args: []interface{}{"va", "vb", "vc"},
		},
		{
			Name: "With three aliases",
			Builder: loukoum.
				Insert("table").
				Columns("a", "b", "c").
				Values([]string{"va", "vb", "vc"}).
				Returning(
					loukoum.Column("a").As("alias_a"),
					loukoum.Column("b").As("alias_b"),
					loukoum.Column("c").As("alias_c"),
				),
			String: fmt.Sprint(
				"INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') ",
				"RETURNING a AS alias_a, b AS alias_b, c AS alias_c",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (a, b, c) VALUES ($1, $2, $3) ",
				"RETURNING a AS alias_a, b AS alias_b, c AS alias_c",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) ",
				"RETURNING a AS alias_a, b AS alias_b, c AS alias_c",
			),
			Args: []interface{}{"va", "vb", "vc"},
		},
	})

	// TODO: expression
}

func TestInsert_Valuer(t *testing.T) {
	when, err := time.Parse(time.RFC3339, "2017-11-23T17:47:27+01:00")
	if err != nil {
		t.Fatal(err)
	}
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "pq.NullTime not null",
			Builder: loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, pq.NullTime{Time: when, Valid: true}),
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', ",
				"true, '2017-11-23 16:47:27+00')",
			),
			Query:      "INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, $3)",
			NamedQuery: "INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, :arg_3)",
			Args:       []interface{}{"tech@ulule.com", true, pq.NullTime{Time: when, Valid: true}},
		},
		{
			Name: "pq.NullTime null",
			Builder: loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, pq.NullTime{}),
			String: fmt.Sprint(
				"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', ",
				"true, NULL)",
			),
			Query:      "INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, $3)",
			NamedQuery: "INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, :arg_3)",
			Args:       []interface{}{"tech@ulule.com", true, pq.NullTime{}},
		},
		{
			Name: "sql.NullString not null",
			Builder: loukoum.
				Insert("table").
				Columns("email", "comment").
				Values("tech@ulule.com", sql.NullString{String: "foobar", Valid: true}),
			String:     "INSERT INTO table (email, comment) VALUES ('tech@ulule.com', 'foobar')",
			Query:      "INSERT INTO table (email, comment) VALUES ($1, $2)",
			NamedQuery: "INSERT INTO table (email, comment) VALUES (:arg_1, :arg_2)",
			Args:       []interface{}{"tech@ulule.com", sql.NullString{String: "foobar", Valid: true}},
		},
		{
			Name: "sql.NullString null",
			Builder: loukoum.
				Insert("table").
				Columns("email", "comment").
				Values("tech@ulule.com", sql.NullString{}),
			String:     "INSERT INTO table (email, comment) VALUES ('tech@ulule.com', NULL)",
			Query:      "INSERT INTO table (email, comment) VALUES ($1, $2)",
			NamedQuery: "INSERT INTO table (email, comment) VALUES (:arg_1, :arg_2)",
			Args:       []interface{}{"tech@ulule.com", sql.NullString{}},
		},
		{
			Name: "sql.NullInt64 not null",
			Builder: loukoum.
				Insert("table").
				Columns("email", "login").
				Values("tech@ulule.com", sql.NullInt64{Int64: 30, Valid: true}),
			String:     "INSERT INTO table (email, login) VALUES ('tech@ulule.com', 30)",
			Query:      "INSERT INTO table (email, login) VALUES ($1, $2)",
			NamedQuery: "INSERT INTO table (email, login) VALUES (:arg_1, :arg_2)",
			Args:       []interface{}{"tech@ulule.com", sql.NullInt64{Int64: 30, Valid: true}},
		},
		{
			Name: "sql.NullInt64 null",
			Builder: loukoum.
				Insert("table").
				Columns("email", "login").
				Values("tech@ulule.com", sql.NullInt64{}),
			String:     "INSERT INTO table (email, login) VALUES ('tech@ulule.com', NULL)",
			Query:      "INSERT INTO table (email, login) VALUES ($1, $2)",
			NamedQuery: "INSERT INTO table (email, login) VALUES (:arg_1, :arg_2)",
			Args:       []interface{}{"tech@ulule.com", sql.NullInt64{}},
		},
		{
			Name: "nil valuer",
			Builder: loukoum.
				Insert("table").
				Columns("email", "login").
				Values("tech@ulule.com", (*valuer)(nil)),
			String:     "INSERT INTO table (email, login) VALUES ('tech@ulule.com', NULL)",
			Query:      "INSERT INTO table (email, login) VALUES ($1, $2)",
			NamedQuery: "INSERT INTO table (email, login) VALUES (:arg_1, :arg_2)",
			Args:       []interface{}{"tech@ulule.com", (*valuer)(nil)},
		},
	})
}

type valuer struct{}

func (valuer) Value() (driver.Value, error) {
	return nil, nil
}
func (*valuer) Scan(src interface{}) error {
	return nil
}

func TestInsert_Set(t *testing.T) {
	RunBuilderTests(t, []BuilderTest{
		{
			Name: "Variadic",
			Builders: []builder.Builder{
				loukoum.Insert("table").Set(
					loukoum.Pair("email", "tech@ulule.com"),
					loukoum.Pair("enabled", true),
					loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				),
				loukoum.Insert("table").Set(
					loukoum.Map{"email": "tech@ulule.com", "enabled": true},
					loukoum.Map{"created_at": loukoum.Raw("NOW()")},
				),
				loukoum.Insert("table").Set(
					map[string]interface{}{"email": "tech@ulule.com"},
					map[string]interface{}{"enabled": true, "created_at": loukoum.Raw("NOW()")},
				),
			},
			String: fmt.Sprint(
				"INSERT INTO table (created_at, email, enabled) ",
				"VALUES (NOW(), 'tech@ulule.com', true)",
			),
			Query: fmt.Sprint(
				"INSERT INTO table (created_at, email, enabled) ",
				"VALUES (NOW(), $1, $2)",
			),
			NamedQuery: fmt.Sprint(
				"INSERT INTO table (created_at, email, enabled) ",
				"VALUES (NOW(), :arg_1, :arg_2)",
			),
			Args: []interface{}{"tech@ulule.com", true},
		},
	})
}
