// Package builder receives user input and generates an AST using "stmt" package.
//
// There is four builder to manipulate an AST: Select, Insert, Update and Delete.
//
// When the AST is ready, you can use String(), NamedQuery() or Query() to generate the underlying query.
// However, be vigilant with String(): it's mainly used for debugging because it's completely vulnerable
// to SQL injection...
package builder
