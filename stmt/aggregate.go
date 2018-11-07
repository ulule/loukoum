package stmt

import (
	"github.com/ulule/loukoum/v2/token"
	"github.com/ulule/loukoum/v2/types"
)

// Count is a aggregate expression.
type Count struct {
	Value      Raw
	IsDistinct bool
	Alias      string
}

// NewCount returns a new Count instance.
func NewCount(value string) Count {
	return Count{
		Value: NewRaw(value),
	}
}

// As is used to give an alias name to the COUNT function.
func (count Count) As(alias string) Count {
	count.Alias = alias
	return count
}

// Distinct is used to define if count has a distinct clause.
func (count Count) Distinct(value bool) Count {
	count.IsDistinct = value
	return count
}

// Write exposes statement as a SQL query.
func (count Count) Write(ctx types.Context) {
	ctx.Write(token.Count.String())
	ctx.Write("(")
	if count.IsDistinct {
		ctx.Write(token.Distinct.String())
		ctx.Write(" ")
	}
	count.Value.Write(ctx)
	ctx.Write(")")
	if count.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(count.Alias)
	}
}

// IsEmpty returns true if statement is undefined.
func (count Count) IsEmpty() bool {
	return count.Value.IsEmpty()
}

func (Count) selectExpression() {}

// Ensure that Count is an SelectExpression
var _ SelectExpression = Count{}

// Max is a aggregate expression.
type Max struct {
	Value Raw
	Alias string
}

// NewMax returns a new Max instance.
func NewMax(value string) Max {
	return Max{
		Value: NewRaw(value),
	}
}

// As is used to give an alias name to the MAX function.
func (max Max) As(alias string) Max {
	max.Alias = alias
	return max
}

// Write exposes statement as a SQL query.
func (max Max) Write(ctx types.Context) {
	ctx.Write(token.Max.String())
	ctx.Write("(")
	max.Value.Write(ctx)
	ctx.Write(")")
	if max.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(max.Alias)
	}
}

// IsEmpty returns true if statement is undefined.
func (max Max) IsEmpty() bool {
	return max.Value.IsEmpty()
}

func (Max) selectExpression() {}

// Ensure that Max is an SelectExpression
var _ SelectExpression = Max{}

// Min is a aggregate expression.
type Min struct {
	Value Raw
	Alias string
}

// NewMin returns a new Min instance.
func NewMin(value string) Min {
	return Min{
		Value: NewRaw(value),
	}
}

// As is used to give an alias name to the MIN function.
func (min Min) As(alias string) Min {
	min.Alias = alias
	return min
}

// Write exposes statement as a SQL query.
func (min Min) Write(ctx types.Context) {
	ctx.Write(token.Min.String())
	ctx.Write("(")
	min.Value.Write(ctx)
	ctx.Write(")")
	if min.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(min.Alias)
	}
}

// IsEmpty returns true if statement is undefined.
func (min Min) IsEmpty() bool {
	return min.Value.IsEmpty()
}

func (Min) selectExpression() {}

// Ensure that Min is an SelectExpression
var _ SelectExpression = Min{}

// Sum is a aggregate expression.
type Sum struct {
	Value Raw
	Alias string
}

// NewSum returns a new Sum instance.
func NewSum(value string) Sum {
	return Sum{
		Value: NewRaw(value),
	}
}

// As is used to give an alias name to the SUM function.
func (sum Sum) As(alias string) Sum {
	sum.Alias = alias
	return sum
}

// Write exposes statement as a SQL query.
func (sum Sum) Write(ctx types.Context) {
	ctx.Write(token.Sum.String())
	ctx.Write("(")
	sum.Value.Write(ctx)
	ctx.Write(")")
	if sum.Alias != "" {
		ctx.Write(" ")
		ctx.Write(token.As.String())
		ctx.Write(" ")
		ctx.Write(sum.Alias)
	}
}

// IsEmpty returns true if statement is undefined.
func (sum Sum) IsEmpty() bool {
	return sum.Value.IsEmpty()
}

func (Sum) selectExpression() {}

// Ensure that Sum is an SelectExpression
var _ SelectExpression = Sum{}
