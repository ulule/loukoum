package stmt

import (
	"github.com/ulule/loukoum/v3/token"
	"github.com/ulule/loukoum/v3/types"
)

// Comment is a comment expression.
type Comment struct {
	Comment string
}

// NewComment returns a new Comment instance.
func NewComment(comment string) Comment {
	return Comment{
		Comment: comment,
	}
}

// Write exposes statement as a SQL query.
func (comment Comment) Write(ctx types.Context) {
	if comment.IsEmpty() {
		return
	}
	ctx.Write(token.Comment.String())
	ctx.Write(" ")
	ctx.Write(comment.Comment)
}

// IsEmpty returns true if statement is undefined.
func (comment Comment) IsEmpty() bool {
	return comment.Comment == ""
}

// Ensure that Comment is a Statement
var _ Statement = Comment{}
