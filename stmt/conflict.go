package stmt

import (
	"github.com/ulule/loukoum/token"
	"github.com/ulule/loukoum/types"
)

// OnConflict is a ON CONFLICT expression.
type OnConflict struct {
	Target ConflictTarget
	Action ConflictAction
}

// NewOnConflict returns a new OnConflict instance.
func NewOnConflict(target ConflictTarget, action ConflictAction) OnConflict {
	return OnConflict{
		Target: target,
		Action: action,
	}
}

// Write exposes statement as a SQL query.
func (conflict OnConflict) Write(ctx *types.Context) {
	if conflict.IsEmpty() {
		return
	}

	ctx.Write(token.On.String())
	ctx.Write(" ")
	ctx.Write(token.Conflict.String())
	ctx.Write(" ")

	if !conflict.Target.IsEmpty() {
		conflict.Target.Write(ctx)
		ctx.Write(" ")
	}

	conflict.Action.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (conflict OnConflict) IsEmpty() bool {
	return conflict.Action == nil || conflict.Action.IsEmpty()
}

// ConflictTarget is a column identifier.
type ConflictTarget struct {
	Columns []Column
}

// NewConflictTarget returns a new ConflictTarget instance.
func NewConflictTarget(columns []Column) ConflictTarget {
	return ConflictTarget{
		Columns: columns,
	}
}

// Write exposes statement as a SQL query.
func (target ConflictTarget) Write(ctx *types.Context) {
	if target.IsEmpty() {
		return
	}

	ctx.Write("(")
	for i := range target.Columns {
		if i != 0 {
			ctx.Write(", ")
		}
		target.Columns[i].Write(ctx)
	}
	ctx.Write(")")
}

// IsEmpty returns true if statement is undefined.
func (target ConflictTarget) IsEmpty() bool {
	return len(target.Columns) == 0
}

// ConflictAction is a action used by ON CONFLICT expression.
// It can be either DO NOTHING, or a DO UPDATE clause.
type ConflictAction interface {
	Statement
	conflictAction()
}

// ConflictUpdateAction is a DO UPDATE clause on ON CONFLICT expression.
type ConflictUpdateAction struct {
	Set Set
}

// NewConflictUpdateAction returns a new ConflictUpdateAction instance.
func NewConflictUpdateAction(set Set) ConflictUpdateAction {
	return ConflictUpdateAction{
		Set: set,
	}
}

// Write exposes statement as a SQL query.
func (action ConflictUpdateAction) Write(ctx *types.Context) {
	ctx.Write(token.Do.String())
	ctx.Write(" ")
	ctx.Write(token.Update.String())
	ctx.Write(" ")
	action.Set.Write(ctx)
}

// IsEmpty returns true if statement is undefined.
func (action ConflictUpdateAction) IsEmpty() bool {
	return action.Set.IsEmpty()
}

func (ConflictUpdateAction) conflictAction() {}

// ConflictNoAction is a DO NOTHING clause on ON CONFLICT expression.
type ConflictNoAction struct{}

// NewConflictNoAction returns a new ConflictNoAction instance.
func NewConflictNoAction() ConflictNoAction {
	return ConflictNoAction{}
}

// Write exposes statement as a SQL query.
func (ConflictNoAction) Write(ctx *types.Context) {
	ctx.Write(token.Do.String())
	ctx.Write(" ")
	ctx.Write(token.Nothing.String())
}

// IsEmpty returns true if statement is undefined.
func (ConflictNoAction) IsEmpty() bool {
	return false
}

func (ConflictNoAction) conflictAction() {}

// Ensure that OnConflict is a Statement
var _ Statement = OnConflict{}

// Ensure that ConflictTarget is a Statement
var _ Statement = ConflictTarget{}

// Ensure that ConflictUpdateAction is a ConflictAction
var _ ConflictAction = ConflictUpdateAction{}

// Ensure that ConflictNoAction is a ConflictAction
var _ ConflictAction = ConflictNoAction{}
