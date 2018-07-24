package parser

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ulule/loukoum/lexer"
	"github.com/ulule/loukoum/token"
)

// ErrAnalyzer is returned when analyze was incomplete.
var ErrAnalyzer = fmt.Errorf("analyze was incomplete")

// AnalyzerOption defines what should be analyzed while scanning statements.
type AnalyzerOption struct {
	// Operation inspects what operation (or command) is performed.
	// Example: INSERT, UPDATE, SELECT, DELETE, etc...
	Operation bool
	// Table inspects what table will handle the given statement.
	Table bool
}

// Continue determines if we should keep scanning statements, or if we have everything we need.
func (o AnalyzerOption) Continue(result *AnalyzerResult) bool {
	if o.Operation && result.Operation == "" {
		return true
	}
	if o.Table && result.Table == "" {
		return true
	}
	return false
}

// AnalyzerResult is the result produces by an analyzer with the given options.
type AnalyzerResult struct {
	// Operation defines the operation (or command).
	// Example: INSERT, UPDATE, SELECT, DELETE, etc...
	Operation string
	// Table defines the table of given statement.
	Table string
}

// Analyze will analyzes given query with options.
func Analyze(query string, option AnalyzerOption) (*AnalyzerResult, error) { // nolint: gocyclo
	lexer := lexer.New(strings.NewReader(query))
	it := lexer.Iterator()

	mode := token.Illegal
	result := &AnalyzerResult{}

	if !option.Continue(result) {
		return result, nil
	}

	for it.HasNext() {
		e := it.Next()
		switch e.Type {
		case token.Select:
			if mode == token.Illegal {
				mode = e.Type
			}
			if !option.Operation {
				continue
			}

			result.Operation = e.Type.String()
			if !option.Continue(result) {
				return result, nil
			}

		case token.Update:
			if mode == token.Illegal {
				mode = e.Type
			}
			if !option.Operation && !option.Table {
				continue
			}

			result.Operation = e.Type.String()
			if !option.Continue(result) {
				return result, nil
			}

			if it.Is(token.Only) {
				it.Next()
			}
			e = it.Next()
			if e.Type != token.Literal {
				return onAnalyzeError(query)
			}
			result.Table = e.Value
			if !option.Continue(result) {
				return result, nil
			}

		case token.Delete:
			if mode == token.Illegal {
				mode = e.Type
			}
			if !option.Operation {
				continue
			}

			result.Operation = e.Type.String()
			if !option.Continue(result) {
				return result, nil
			}

		case token.Insert:
			if mode == token.Illegal {
				mode = e.Type
			}
			if !option.Operation && !option.Table {
				continue
			}

			result.Operation = e.Type.String()
			if !option.Continue(result) {
				return result, nil
			}

			e = it.Next()
			if e.Type != token.Into {
				return onAnalyzeError(query)
			}
			e = it.Next()
			if e.Type != token.Literal {
				return onAnalyzeError(query)
			}
			result.Table = e.Value
			if !option.Continue(result) {
				return result, nil
			}

		case token.From:
			if !option.Table {
				continue
			}

			if mode == token.Select || mode == token.Delete {
				if it.Is(token.Only) {
					it.Next()
				}
				e = it.Next()
				if e.Type != token.Literal {
					return onAnalyzeError(query)
				}
				result.Table = e.Value
				if !option.Continue(result) {
					return result, nil
				}
			}

		case token.With:
			if it.Is(token.Recursive) {
				it.Next()
			}
			if !onAnalyzeWithQuery(it, result) {
				return onAnalyzeError(query)
			}
			for it.Is(token.Comma) {
				it.Next()
				if !onAnalyzeWithQuery(it, result) {
					return onAnalyzeError(query)
				}
			}
		}
	}

	if option.Continue(result) {
		return onAnalyzeError(query)
	}

	return result, nil
}

func onAnalyzeError(query string) (*AnalyzerResult, error) {
	return nil, errors.Wrapf(ErrAnalyzer, "parsing error with: %s", query)
}

// onAnalyzeWithQuery will consumes current "WITH" query in iteratee.
func onAnalyzeWithQuery(it *lexer.Iteratee, result *AnalyzerResult) bool {
	e := it.Next()
	if e.Type != token.Literal {
		return false
	}
	for it.Is(token.Literal) {
		it.Next()
		if it.Is(token.Comma) {
			it.Next()
		}
	}
	e = it.Next()
	if e.Type != token.As {
		return false
	}
	e = it.Next()
	if e.Type != token.LParen {
		return false
	}
	level := 1
	for it.HasNext() {
		e = it.Next()
		if e.Type == token.LParen {
			level++
		}
		if e.Type == token.RParen {
			level--
		}
		if level == 0 {
			return true
		}
	}
	return level == 0
}
