package logging

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLogger(t *testing.T) {
	t.Run("init log", func(t *testing.T) {
		err := InitLog(-1, "")
		require.NoError(t, err)

		ctx := context.Background()
		logger := FromContext(ctx)

		require.Equal(t, fallbackLogger, logger)

		ctx = WithLogger(ctx, logger)
		require.Equal(t, logger, FromContext(ctx))
	})
}
