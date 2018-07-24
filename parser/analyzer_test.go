package parser_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/parser"
)

func TestAnalyze(t *testing.T) {
	is := require.New(t)

	scenarios := []struct {
		Query    string
		Option   parser.AnalyzerOption
		Expected parser.AnalyzerResult
	}{
		{
			// Scenario #1
			Query: "SELECT a, b, c FROM foobar",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "SELECT",
				Table:     "foobar",
			},
		},
		{
			// Scenario #2
			Query: "SELECT a, b, c FROM foobar AS example",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "SELECT",
				Table:     "foobar",
			},
		},
		{
			// Scenario #3
			Query: "SELECT a, b, c FROM test1 INNER JOIN test2 ON test1.id = test2.fk_id",
			Option: parser.AnalyzerOption{
				Operation: false,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "",
				Table:     "test1",
			},
		},
		{
			// Scenario #4
			Query: "SELECT * FROM table WHERE id = 596 AND enabled = true",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "SELECT",
				Table:     "table",
			},
		},
		{
			// Scenario #5
			Query: "INSERT INTO table (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW())",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "INSERT",
				Table:     "table",
			},
		},
		{
			// Scenario #6
			Query: fmt.Sprint(
				"INSERT INTO table AS foobar (email, enabled, created_at) VALUES ('tech@ulule.com', true, NOW()) ",
				"ON CONFLICT (email) DO UPDATE SET created_at = NOW(), enabled = true",
			),
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "INSERT",
				Table:     "table",
			},
		},
		{
			// Scenario #7
			Query: "UPDATE table SET a = 1, foo = 2",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "UPDATE",
				Table:     "table",
			},
		},
		{
			// Scenario #8
			Query: "UPDATE ONLY table SET a = 1, foo = 2",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "UPDATE",
				Table:     "table",
			},
		},
		{
			// Scenario #9
			Query: "DELETE FROM table USING foobar",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "DELETE",
				Table:     "table",
			},
		},
		{
			// Scenario #10
			Query: "DELETE FROM table WHERE id = 1",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "DELETE",
				Table:     "table",
			},
		},
		{
			// Scenario #11
			Query: "DELETE ONLY FROM table WHERE id = 1",
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "DELETE",
				Table:     "table",
			},
		},
		{
			// Scenario #12
			Query: fmt.Sprint(
				"WITH contributors AS (SELECT DISTINCT user_id FROM contribution WHERE (deleted_at IS NULL)) ",
				"UPDATE users SET newsletter_subscribed = true FROM contributors ",
				"WHERE (users.id = contributors.user_id)",
			),
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "UPDATE",
				Table:     "users",
			},
		},
		{
			// Scenario #13
			Query: fmt.Sprint(
				"WITH contributors AS (SELECT DISTINCT user_id FROM contribution WHERE (deleted_at IS NULL)), ",
				"commentators AS (SELECT DISTINCT user_id FROM comment WHERE (deleted_at IS NULL)) ",
				"UPDATE users SET newsletter_subscribed = true ",
				"WHERE ((users.id IN (SELECT user_id FROM contributors)) ",
				"OR (users.id IN (SELECT user_id FROM commentators)))",
			),
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "UPDATE",
				Table:     "users",
			},
		},
		{
			// Scenario #14
			Query: fmt.Sprint(
				"WITH withdrawn AS (SELECT SUM(amount) AS amount, project_id FROM withdrawal GROUP BY project_id), ",
				"contributions AS (SELECT COUNT(*) AS count, project_id FROM contribution GROUP BY project_id) ",
				"SELECT SUM(project.amount_raised - withdrawn.amount) FROM project ",
				"LEFT JOIN withdrawn ON withdrawn.project_id = project.id ",
				"LEFT JOIN contributions ON contributions.project_id = project.id ",
				"WHERE ((((project.amount_raised > 0) AND (contributions.count > 10)) AND ",
				"(project.deleted_at IS NULL)) AND (project.amount_raised > withdrawn.amount))",
			),
			Option: parser.AnalyzerOption{
				Operation: true,
				Table:     true,
			},
			Expected: parser.AnalyzerResult{
				Operation: "SELECT",
				Table:     "project",
			},
		},
	}

	for i, scenario := range scenarios {
		message := fmt.Sprintf("scenario #%d", (i + 1))
		result, err := parser.Analyze(scenario.Query, scenario.Option)
		is.NoError(err, message)
		is.NotNil(result, message)
		is.Equal(scenario.Expected, *result, message)
	}
}
