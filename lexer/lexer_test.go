package lexer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/v3/lexer"
	"github.com/ulule/loukoum/v3/token"
)

type LexScenario struct {
	Input  string
	Tokens []token.Token
}

func execute(t *testing.T, scenarios []LexScenario) {
	is := require.New(t)

	for i, scenario := range scenarios {

		l := lexer.New(strings.NewReader(scenario.Input))

		for y, expected := range scenario.Tokens {
			message := fmt.Sprintf("Scenario #%d / Token #%d", (i + 1), (y + 1))

			actual := l.Next()

			is.Equal(expected.Type, actual.Type, message)
			is.Equal(expected.Value, actual.Value, message)
		}

		actual := l.Next()
		is.Equal(token.EOF, actual.Type, "EOF was expected")
	}
}

func TestNextToken(t *testing.T) {

	tests := []LexScenario{}

	// Scenario #1: Check EOF on empty source
	tests = append(tests, LexScenario{
		Input: ``,
		Tokens: []token.Token{
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
			token.New(token.EOF, ""),
		},
	})

	// Scenario #2: A simple SELECT query
	tests = append(tests, LexScenario{
		Input: `
			SELECT * FROM foobar WHERE id = 2;
		`,
		Tokens: []token.Token{
			token.New(token.Select, "SELECT"),
			token.New(token.Asterisk, "*"),
			token.New(token.From, "FROM"),
			token.New(token.Literal, "foobar"),
			token.New(token.Where, "WHERE"),
			token.New(token.Literal, "id"),
			token.New(token.Equals, "="),
			token.New(token.Literal, "2"),
			token.New(token.Semicolon, ";"),
		},
	})

	// Scenario #3: Another simple SELECT query with multiples columns
	tests = append(tests, LexScenario{
		Input: `
			SELECT a,b,c FROM foobar;
		`,
		Tokens: []token.Token{
			token.New(token.Select, "SELECT"),
			token.New(token.Literal, "a"),
			token.New(token.Comma, ","),
			token.New(token.Literal, "b"),
			token.New(token.Comma, ","),
			token.New(token.Literal, "c"),
			token.New(token.From, "FROM"),
			token.New(token.Literal, "foobar"),
			token.New(token.Semicolon, ";"),
		},
	})

	// Scenario #4: Detect if newline is ignored
	tests = append(tests, LexScenario{
		Input: `
			SELECT a, b, c FROM foobar
		`,
		Tokens: []token.Token{
			token.New(token.Select, "SELECT"),
			token.New(token.Literal, "a"),
			token.New(token.Comma, ","),
			token.New(token.Literal, "b"),
			token.New(token.Comma, ","),
			token.New(token.Literal, "c"),
			token.New(token.From, "FROM"),
			token.New(token.Literal, "foobar"),
		},
	})

	// Scenario #5: Detect if EOF is ignored
	tests = append(tests, LexScenario{
		Input: `SELECT a, b, c FROM foobar`,
		Tokens: []token.Token{
			token.New(token.Select, "SELECT"),
			token.New(token.Literal, "a"),
			token.New(token.Comma, ","),
			token.New(token.Literal, "b"),
			token.New(token.Comma, ","),
			token.New(token.Literal, "c"),
			token.New(token.From, "FROM"),
			token.New(token.Literal, "foobar"),
		},
	})

	// Scenario #6: A subquery using an INNER JOIN
	tests = append(tests, LexScenario{
		Input: `INNER JOIN test2 ON test2.id = test.fk_id`,
		Tokens: []token.Token{
			token.New(token.Inner, "INNER"),
			token.New(token.Join, "JOIN"),
			token.New(token.Literal, "test2"),
			token.New(token.On, "ON"),
			token.New(token.Literal, "test2.id"),
			token.New(token.Equals, "="),
			token.New(token.Literal, "test.fk_id"),
		},
	})

	// Scenario #7: A simple delete query
	tests = append(tests, LexScenario{
		Input: `DELETE FROM test2 WHERE id = 5`,
		Tokens: []token.Token{
			token.New(token.Delete, "DELETE"),
			token.New(token.From, "FROM"),
			token.New(token.Literal, "test2"),
			token.New(token.Where, "WHERE"),
			token.New(token.Literal, "id"),
			token.New(token.Equals, "="),
			token.New(token.Literal, "5"),
		},
	})

	execute(t, tests)
}
