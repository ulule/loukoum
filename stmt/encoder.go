package stmt

import (
	"time"
)

// StatementEncoder can encode a value as a statement to creates a Expression instance.
type StatementEncoder interface {
	Statement() Statement
}

// StringEncoder can encode a value as a string to creates a Value instance.
type StringEncoder interface {
	String() string
}

// Int64Encoder can encode a value as a int64 to creates a Value instance.
type Int64Encoder interface {
	Int64() int64
}

// BoolEncoder can encode a value as a bool to creates a Value instance.
type BoolEncoder interface {
	Bool() bool
}

// TimeEncoder can encode a value as a time.Time to creates a Value instance.
type TimeEncoder interface {
	Time() time.Time
}
