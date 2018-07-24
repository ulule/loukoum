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
func (o AnalyzerOption) Continue(result AnalyzerResult) bool {
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
func Analyze(query string, option AnalyzerOption) (AnalyzerResult, error) { // nolint: gocyclo
	lexer := lexer.New(strings.NewReader(query))
	it := lexer.Iterator()

	mode := token.Illegal
	result := AnalyzerResult{}

	if !option.Continue(result) {
		return result, nil
	}

	// NOTE (novln): Not working with "WITH" statements...

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

		}
	}

	if option.Continue(result) {
		return onAnalyzeError(query)
	}

	return result, nil
}

func onAnalyzeError(query string) (AnalyzerResult, error) {
	return AnalyzerResult{}, errors.Wrapf(ErrAnalyzer, "parsing error with: %s", query)
}
