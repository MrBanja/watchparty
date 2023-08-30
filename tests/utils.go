package tests

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type Receiver[Res any] interface {
	Receive() (*Res, error)
}

func mustReceive[Res any](t *testing.T, bidi Receiver[Res]) *Res {
	res, err := bidi.Receive()
	require.NoError(t, err)
	require.NotNil(t, res)
	return res
}
