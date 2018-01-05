package types

import (
	"time"
)

// StringEncoder can encode a value as a string to create a stmt.Value instance.
type StringEncoder interface {
	String() string
}

// Int64Encoder can encode a value as a int64 to create a stmt.Value instance.
type Int64Encoder interface {
	Int64() int64
}

// BoolEncoder can encode a value as a bool to create a stmt.Value instance.
type BoolEncoder interface {
	Bool() bool
}

// TimeEncoder can encode a value as a time.Time to create a stmt.Value instance.
type TimeEncoder interface {
	Time() time.Time
}
