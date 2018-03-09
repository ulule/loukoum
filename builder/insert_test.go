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

var inserttests = []BuilderTest{
	{
		Name:       "Columns With columns",
		Builder:    loukoum.Insert("table").Columns("a", "b", "c"),
		String:     "INSERT INTO table (a, b, c)",
		Query:      "INSERT INTO table (a, b, c)",
		NamedQuery: "INSERT INTO table (a, b, c)",
	},
	{
		Name:       "Columns Without columns",
		Builder:    loukoum.Insert("table"),
		String:     "INSERT INTO table",
		Query:      "INSERT INTO table",
		NamedQuery: "INSERT INTO table",
	},
	{
		Name: "Values With columns",
		Builder: loukoum.Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}),
		String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc')",
		Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3)",
		NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3)",
		Args:       []interface{}{"va", "vb", "vc"},
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Values Without columns",
		Builder: loukoum.
			Insert("table").
			Values([]string{"va", "vb", "vc"}),
		String:     "INSERT INTO table VALUES ('va', 'vb', 'vc')",
		Query:      "INSERT INTO table VALUES ($1, $2, $3)",
		NamedQuery: "INSERT INTO table VALUES (:arg_1, :arg_2, :arg_3)",
		Args:       []interface{}{"va", "vb", "vc"},
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "OnConflict Do nothing without target",
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "OnConflict Do nothing A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", loukoum.DoNothing()),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "OnConflict Do nothing B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict(loukoum.Column("email"), loukoum.DoNothing()),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "OnConflict Do nothing C",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", "uuid", loukoum.DoNothing()),
		String: fmt.Sprint(
			"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
			"ON CONFLICT (email, uuid) DO NOTHING",
		),
		Query: fmt.Sprint(
			"INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NOW()) ",
			"ON CONFLICT (email, uuid) DO NOTHING",
		),
		NamedQuery: fmt.Sprint(
			"INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NOW()) ",
			"ON CONFLICT (email, uuid) DO NOTHING",
		),
		Args: []interface{}{"tech@ulule.com", true},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "OnConflict Do nothing D",
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "OnConflict Do update A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict("email", loukoum.DoUpdate(
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				loukoum.Pair("enabled", true),
			)),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
			"arg_3": true,
		},
	},
	{
		Name: "OnConflict Do update B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
			OnConflict(loukoum.Column("email"), loukoum.DoUpdate(
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
				loukoum.Pair("enabled", true),
			)),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
			"arg_3": true,
		},
	},
	{
		Name: "OnConflict Do update C",
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
			"arg_3": true,
		},
	},
	{
		Name: "OnConflict Corner case A",
		Failure: func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict()
		},
	},
	{
		Name: "OnConflict Corner case B",
		Failure: func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email")
		},
	},
	{
		Name: "OnConflict Corner case C",
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
		Name: "OnConflict Corner case D",
		Failure: func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", 6700)
		},
	},
	{
		Name: "OnConflict Corner case E",
		Failure: func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict(569)
		},
	},
	{
		Name: "OnConflict Corner case F",
		Failure: func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict("email", "uuid")
		},
	},
	{
		Name: "OnConflict Corner case G",
		Failure: func() builder.Builder {
			return loukoum.
				Insert("table").
				Columns("email", "enabled", "created_at").
				Values("tech@ulule.com", true, loukoum.Raw("NOW()")).
				OnConflict(loukoum.Column("email"), loukoum.Column("uuid"), loukoum.Column("reference"))
		},
	},
	{
		Name: "Returning One column",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a"),
		String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a",
		Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a",
		NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a",
		Args:       []interface{}{"va", "vb", "vc"},
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Returning Many columns A",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a", "b"),
		String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a, b",
		Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a, b",
		NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a, b",
		Args:       []interface{}{"va", "vb", "vc"},
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Returning Many columns B",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a", "b", "c"),
		String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a, b, c",
		Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a, b, c",
		NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a, b, c",
		Args:       []interface{}{"va", "vb", "vc"},
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Returning With aliases A",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning(loukoum.Column("a").As("alias_a")),
		String:     "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a AS alias_a",
		Query:      "INSERT INTO table (a, b, c) VALUES ($1, $2, $3) RETURNING a AS alias_a",
		NamedQuery: "INSERT INTO table (a, b, c) VALUES (:arg_1, :arg_2, :arg_3) RETURNING a AS alias_a",
		Args:       []interface{}{"va", "vb", "vc"},
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Returning With aliases B",
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
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Returning With aliases C",
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
		NamedArgs: map[string]interface{}{
			"arg_1": "va",
			"arg_2": "vb",
			"arg_3": "vc",
		},
	},
	{
		Name: "Valuer pq.NullTime A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, pq.NullTime{Time: now, Valid: true}),
		String: fmt.Sprint(
			"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', ",
			"true, ", format.Time(now), ")",
		),
		Query:      "INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, $3)",
		NamedQuery: "INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, :arg_3)",
		Args:       []interface{}{"tech@ulule.com", true, now},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
			"arg_3": now,
		},
	},
	{
		Name: "Valuer pq.NullTime B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "enabled", "created_at").
			Values("tech@ulule.com", true, pq.NullTime{}),
		String: fmt.Sprint(
			"INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', ",
			"true, NULL)",
		),
		Query:      "INSERT INTO table (email, enabled, created_at) VALUES ($1, $2, NULL)",
		NamedQuery: "INSERT INTO table (email, enabled, created_at) VALUES (:arg_1, :arg_2, NULL)",
		Args:       []interface{}{"tech@ulule.com", true},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "Valuer sql.NullString A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "comment").
			Values("tech@ulule.com", sql.NullString{String: "foobar", Valid: true}),
		String:     "INSERT INTO table (email, comment) VALUES ('tech@ulule.com', 'foobar')",
		Query:      "INSERT INTO table (email, comment) VALUES ($1, $2)",
		NamedQuery: "INSERT INTO table (email, comment) VALUES (:arg_1, :arg_2)",
		Args:       []interface{}{"tech@ulule.com", "foobar"},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": "foobar",
		},
	},
	{
		Name: "Valuer sql.NullString B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "comment").
			Values("tech@ulule.com", sql.NullString{}),
		String:     "INSERT INTO table (email, comment) VALUES ('tech@ulule.com', NULL)",
		Query:      "INSERT INTO table (email, comment) VALUES ($1, NULL)",
		NamedQuery: "INSERT INTO table (email, comment) VALUES (:arg_1, NULL)",
		Args:       []interface{}{"tech@ulule.com"},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
		},
	},
	{
		Name: "Valuer sql.NullInt64 A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "login").
			Values("tech@ulule.com", sql.NullInt64{Int64: 30, Valid: true}),
		String:     "INSERT INTO table (email, login) VALUES ('tech@ulule.com', 30)",
		Query:      "INSERT INTO table (email, login) VALUES ($1, $2)",
		NamedQuery: "INSERT INTO table (email, login) VALUES (:arg_1, :arg_2)",
		Args:       []interface{}{"tech@ulule.com", int64(30)},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": int64(30),
		},
	},
	{
		Name: "Valuer sql.NullInt64 B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "login").
			Values("tech@ulule.com", sql.NullInt64{}),
		String:     "INSERT INTO table (email, login) VALUES ('tech@ulule.com', NULL)",
		Query:      "INSERT INTO table (email, login) VALUES ($1, NULL)",
		NamedQuery: "INSERT INTO table (email, login) VALUES (:arg_1, NULL)",
		Args:       []interface{}{"tech@ulule.com"},
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
		},
	},
	{
		Name: "Set Variadic with Pair type",
		Builder: loukoum.
			Insert("table").
			Set(
				loukoum.Pair("email", "tech@ulule.com"),
				loukoum.Pair("enabled", true),
				loukoum.Pair("created_at", loukoum.Raw("NOW()")),
			),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "Set Variadic with Map type",
		Builder: loukoum.
			Insert("table").
			Set(
				loukoum.Map{"email": "tech@ulule.com", "enabled": true},
				loukoum.Map{"created_at": loukoum.Raw("NOW()")},
			),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
	{
		Name: "Set Variadic with string / interface map",
		Builder: loukoum.
			Insert("table").
			Set(
				map[string]interface{}{"email": "tech@ulule.com"},
				map[string]interface{}{"enabled": true, "created_at": loukoum.Raw("NOW()")},
			),
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
		NamedArgs: map[string]interface{}{
			"arg_1": "tech@ulule.com",
			"arg_2": true,
		},
	},
}

func TestInsert(t *testing.T) {
	for _, tt := range inserttests {
		t.Run(tt.Name, tt.Run)
	}
}
