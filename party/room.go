package party

import (
	"sync"

	gen "github.com/mrbanja/watchparty/protocol/gen-go"

	"github.com/bufbuild/connect-go"
	"go.uber.org/zap"
)

type Room struct {
	Name string

	participantMu sync.RWMutex
	participants  map[*Participant]struct{}

	logger *zap.Logger
}

func NewRoom(name string, logger *zap.Logger) *Room {
	return &Room{
		Name:         name,
		participants: make(map[*Participant]struct{}),
		logger:       logger.Named("room").With(zap.String("Room Name", name)),
	}
}

func (r *Room) AddParticipant(c *connect.BidiStream[gen.RoomRequest, gen.RoomResponse]) *Participant {
	p := newParticipant(c)
	r.participantMu.Lock()
	r.participants[p] = struct{}{}
	r.participantMu.Unlock()
	r.logger.Info("Peer connected to room", zap.String("Addr", p.peer.Addr), zap.String("Participant ID", p.ID))
	return p
}

func (r *Room) RemoveParticipant(p *Participant) {
	defer r.logger.Info("Peer disconnected from the room", zap.String("Addr", p.peer.Addr))
	r.participantMu.Lock()
	delete(r.participants, p)
	r.participantMu.Unlock()
}

func (r *Room) Broadcast(msg *gen.RoomResponse) {
	r.BroadcastExcept(msg, nil)
}

func (r *Room) BroadcastExcept(msg *gen.RoomResponse, participant *Participant) {
	r.participantMu.RLock()
	defer r.participantMu.RUnlock()

	for p := range r.participants {
		if p == participant {
			continue
		}
		if err := p.conn.Send(msg); err != nil {
			r.logger.Error("Error sending to peer", zap.Error(err), zap.String("Peer ID", p.ID))
		}
	}
	r.logger.Info("Broadcast to the room", zap.Any("MSG", msg))
}

func (r *Room) GetParticipantByID(ID string) *Participant {
	r.participantMu.RLock()
	defer r.participantMu.RUnlock()
	for p := range r.participants {
		if p.ID == ID {
			return p
		}
	}
	return nil
}
