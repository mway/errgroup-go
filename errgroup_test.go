package errgroup_test

import (
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mway.dev/errgroup"
	"go.uber.org/multierr"
)

var (
	errA = errors.New("a")
	errB = errors.New("b")
	errC = errors.New("c")
)

func TestErrGroupNoErrors(t *testing.T) {
	var g1 errgroup.Group

	g1.Add(
		func() error { return nil },
		func() error { return nil },
		func() error { return nil },
	)

	require.NoError(t, g1.Wait())
}

func TestErrGroupParallel(t *testing.T) {
	var g errgroup.Group
	g.Add(
		func() error {
			time.Sleep(100 * time.Millisecond)
			return errA
		},
		func() error {
			time.Sleep(50 * time.Millisecond)
			return errB
		},
		func() error {
			return errC
		},
	)
}

func TestErrGroupInline(t *testing.T) {
	var (
		expectErr = multierr.Combine(errA, errB, errC)
		g         = errgroup.New(errgroup.WithInline())
	)

	g.Add(
		func() error {
			time.Sleep(100 * time.Millisecond)
			return errA
		},
		func() error {
			time.Sleep(50 * time.Millisecond)
			return errB
		},
		func() error {
			return errC
		},
	)

	require.EqualError(t, g.Wait(), expectErr.Error())
}

func TestErrGroupFirstOnly(t *testing.T) {
	var (
		expectErr = errA
		g         = errgroup.New(
			errgroup.WithInline(),
			errgroup.WithFirstOnly(),
		)
	)

	g.Add(
		func() error {
			return errA
		},
		func() error {
			return errB
		},
		func() error {
			return errC
		},
	)

	require.EqualError(t, g.Wait(), expectErr.Error())
}

func TestErrGroupIgnoredErrors(t *testing.T) {
	var (
		expectErr = errC
		g         = errgroup.New(
			errgroup.WithInline(),
			errgroup.WithIgnoredErrors(io.EOF),
		)
	)

	g.Add(
		func() error {
			return nil
		},
		func() error {
			return io.EOF
		},
		func() error {
			return errC
		},
	)

	require.EqualError(t, g.Wait(), expectErr.Error())
}
