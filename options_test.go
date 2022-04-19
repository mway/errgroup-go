package errgroup_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errgroup"
)

func TestOptionsWith(t *testing.T) {
	var (
		previous = errgroup.DefaultOptions().With(errgroup.WithFirstOnly())
		updated  = previous.With(errgroup.DefaultOptions().With(errgroup.WithInline()))
	)

	require.True(t, previous.FirstOnly)
	require.False(t, previous.Inline)
	require.False(t, updated.FirstOnly)
	require.True(t, updated.Inline)
}
