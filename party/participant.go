package party

import (
	"github.com/bufbuild/connect-go"
	gen "github.com/mrbanja/watchparty/protocol/gen-go"
)

type Participant struct {
	ID string

	peer connect.Peer
	conn *connect.BidiStream[gen.RoomRequest, gen.RoomResponse]
}

func newParticipant(c *connect.BidiStream[gen.RoomRequest, gen.RoomResponse]) *Participant {
	peerID := c.RequestHeader().Get("X-Client-Id")
	if peerID == "" {
		peerID = "empty"
	}
	return &Participant{peer: c.Peer(), ID: peerID, conn: c}
}
