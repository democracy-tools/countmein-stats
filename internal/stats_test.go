package internal_test

import (
	"testing"

	"github.com/democracy-tools/countmein-stats/internal"
	"github.com/democracy-tools/countmein-stats/internal/bq"
	"github.com/democracy-tools/countmein-stats/internal/gcs"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {

	require.NoError(t, internal.Run(bq.NewInMemoryClient(), gcs.NewInMemoryClient()))
}
