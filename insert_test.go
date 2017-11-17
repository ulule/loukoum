package loukoum_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ulule/loukoum"
)

func TestInsert(t *testing.T) {
	is := require.New(t)

	// With columns
	{
		query := loukoum.
			Insert("table").
			Columns("a", "b", "c").
			Values([]string{"va", "vb", "vc"})

		is.Equal("INSERT INTO table (a, b, c) VALUES (va, vb, vc)", query.String())
	}

	// Without columns
	{
		query := loukoum.
			Insert("table").
			Values([]string{"va", "vb", "vc"})

		is.Equal("INSERT INTO table VALUES (va, vb, vc)", query.String())
	}
}
