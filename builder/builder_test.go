package builder_test

import (
	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum/builder"
)

func Failure(is *require.Assertions, callback func() builder.Builder) {
	is.Panics(func() {
		_ = callback().String()
	})
}
