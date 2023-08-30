package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/themakers/bdd"
)

func TestBadge(t *testing.T) {
	Scenario(t, "Badge", func(t *testing.T, runID string) {
		Test(t, "Get Badge", func() {
			req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8002/room_id/peer_status/participant_id", nil)
			require.NoError(t, err)
			require.NotNil(t, req)

			req.Header.Set("X-Client-Id", "TEST-X-ID")

			resp, err := http.DefaultClient.Do(req)
			require.NotNil(t, resp)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, resp.StatusCode)
		})
	})
}
