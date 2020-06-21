package webserver

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestWebServer(t *testing.T) {
	t.Run("web server", func(t *testing.T) {
		w := NewWebServer(nil, "[::1]:56001")
		go func() {
			<-time.After(time.Second)
			err := w.Shutdown(time.Millisecond * 100)
			require.NoError(t, err)
		}()

		err := w.Start()
		require.NoError(t, err)
	})
}
