package errgroup_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errgroup"
)

func TestOptionsWith(t *testing.T) {
	var (
		previous = errgroup.DefaultOptions().With(
			errgroup.WithFirstOnly(),
			errgroup.WithIgnoredErrors(io.EOF),
		)
		updated = previous.With(errgroup.DefaultOptions().With(errgroup.WithInline()))
	)

	require.True(t, previous.FirstOnly)
	require.False(t, previous.Inline)
	require.Len(t, previous.IgnoredErrors, 1)

	require.False(t, updated.FirstOnly)
	require.True(t, updated.Inline)
	require.Len(t, updated.IgnoredErrors, 1)
}
