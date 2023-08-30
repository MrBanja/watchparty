package tests

import (
	"context"
	"testing"

	gen_go "github.com/mrbanja/watchparty-proto/gen-go"

	"github.com/mrbanja/watchparty-proto/gen-go/protocolconnect"
	"github.com/mrbanja/watchparty/tools/http_client"

	"github.com/stretchr/testify/require"
	. "github.com/themakers/bdd"
)

func TestGame(t *testing.T) {
	hc := http_client.New()
	Scenario(t, "Party", func(t *testing.T, runID string) {
		partyClient := protocolconnect.NewPartyServiceClient(hc, "http://127.0.0.1:8002")

		Act(t, "Party Stream", func() {
			Test(t, "Connect", func() {
				bidi := partyClient.JoinRoom(context.TODO())
				require.NotNil(t, bidi)

				err := bidi.Send(&gen_go.RoomRequest{
					Data: &gen_go.RoomRequest_Connect{
						Connect: &gen_go.Connect{
							RoomName: "Party",
						},
					},
				})
				require.NoError(t, err)

				res := mustReceive[gen_go.RoomResponse](t, bidi)
				t.Log(res)
			})
		})
	})
}
