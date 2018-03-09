package builder_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
	"github.com/ulule/loukoum/builder"
	"github.com/ulule/loukoum/format"
)

var now = time.Now()

type BuilderTest struct {
	Name       string
	Builder    builder.Builder
	String     string
	Query      string
	Args       []interface{}
	NamedQuery string
	NamedArgs  map[string]interface{}
	Failure    func() builder.Builder
}

func (b *BuilderTest) Run(t *testing.T) {
	if b.Failure != nil {
		t.Run("Failure", func(t *testing.T) {
			require.Panics(t, func() {
				_ = b.Failure()
			})
		})
	}
	if b.Builder == nil {
		return
	}
	t.Run("String", func(t *testing.T) {
		if b.String != "" {
			require.Equal(t, b.String, b.Builder.String())
		}
	})
	t.Run("Query", func(t *testing.T) {
		query, args := b.Builder.Query()
		if b.Query != "" {
			require.Equal(t, b.Query, query)
			require.Equal(t, b.Args, args)
		}
	})
	t.Run("NamedQuery", func(t *testing.T) {
		query, args := b.Builder.NamedQuery()
		if b.NamedQuery != "" {
			require.Equal(t, b.NamedQuery, query)
			require.Equal(t, b.NamedArgs, args)
		}
	})
}

var inserttests = []BuilderTest{
	{
		Name:    "Columns With columns",
		Builder: loukoum.Insert("table").Columns("a", "b", "c"),
		String:  "INSERT INTO table (a, b, c)",
	},
	{
		Name:    "Columns Without columns",
		Builder: loukoum.Insert("table"),
		String:  "INSERT INTO table",
	},
	{
		Name: "Values With columns A",
		Builder: loukoum.Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}),
		String: "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc')",
	},
	{
		Name: "Values With Columns B",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values("va", "vb", "vc"),
		String: "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc')",
	},
	{
		Name: "Values Without columns A",
		Builder: loukoum.
			Insert("table").
			Values([]string{"va", "vb", "vc"}),
		String: "INSERT INTO table VALUES ('va', 'vb', 'vc')",
	},
	{
		Name: "Values Without columns B",
		Builder: loukoum.
			Insert("table").
			Values("va", "vb", "vc"),
		String: "INSERT INTO table VALUES ('va', 'vb', 'vc')",
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
		String: "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a",
	},
	{
		Name: "Returning Many columns A",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a", "b"),
		String: "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a, b",
	},
	{
		Name: "Returning Many columns B",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning("a", "b", "c"),
		String: "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a, b, c",
	},
	{
		Name: "Returning With aliases A",
		Builder: loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"}).
			Returning(loukoum.Column("a").As("alias_a")),
		String: "INSERT INTO table (a, b, c) VALUES ('va', 'vb', 'vc') RETURNING a AS alias_a",
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
	},
	{
		Name: "Valuer sql.NullString A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "comment").
			Values("tech@ulule.com", sql.NullString{String: "foobar", Valid: true}),
		String: "INSERT INTO table (email, comment) VALUES ('tech@ulule.com', 'foobar')",
	},
	{
		Name: "Valuer sql.NullString B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "comment").
			Values("tech@ulule.com", sql.NullString{}),
		String: "INSERT INTO table (email, comment) VALUES ('tech@ulule.com', NULL)",
	},
	{
		Name: "Valuer sql.NullInt64 A",
		Builder: loukoum.
			Insert("table").
			Columns("email", "login").
			Values("tech@ulule.com", sql.NullInt64{Int64: 30, Valid: true}),
		String: "INSERT INTO table (email, login) VALUES ('tech@ulule.com', 30)",
	},
	{
		Name: "Valuer sql.NullInt64 B",
		Builder: loukoum.
			Insert("table").
			Columns("email", "login").
			Values("tech@ulule.com", sql.NullInt64{}),
		String: "INSERT INTO table (email, login) VALUES ('tech@ulule.com', NULL)",
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
	},
}

func TestInsert(t *testing.T) {
	for _, tt := range inserttests {
		t.Run(tt.Name, tt.Run)
	}
}
