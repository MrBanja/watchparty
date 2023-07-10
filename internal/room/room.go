package room

import (
	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
	"sync"
)

type Participant struct {
	conn *websocket.Conn
}

func newParticipant(c *websocket.Conn) *Participant {
	return &Participant{conn: c}
}

func (p *Participant) ReadMessage() ([]byte, error) {
	_, msg, err := p.conn.ReadMessage()
	return msg, err
}

func (p *Participant) WriteMessage(msg string) error {
	return p.conn.WriteMessage(1, []byte(msg))
}

type Room struct {
	Name string

	participantMu sync.RWMutex
	participants  map[*Participant]struct{}
}

func New(name string) *Room {
	return &Room{
		Name:         name,
		participants: make(map[*Participant]struct{}),
	}
}

func (r *Room) AddParticipant(conn *websocket.Conn) *Participant {
	defer zap.S().Infof("Peer connected %s to room %s\n", conn.RemoteAddr(), r.Name)
	p := newParticipant(conn)
	r.participantMu.Lock()
	r.participants[p] = struct{}{}
	r.participantMu.Unlock()
	return p
}

func (r *Room) RemoveParticipant(p *Participant) {
	defer zap.S().Infof("Peer disconnected %s from room %s\n", p.conn.RemoteAddr(), r.Name)
	r.participantMu.Lock()
	delete(r.participants, p)
	r.participantMu.Unlock()
}

func (r *Room) Broadcast(msg string) {
	r.BroadcastExcept(msg, nil)
}

func (r *Room) BroadcastExcept(msg string, participant *Participant) {
	r.participantMu.RLock()
	defer r.participantMu.RUnlock()

	for p := range r.participants {
		if p == participant {
			continue
		}
		if err := p.WriteMessage(msg); err != nil {
			zap.S().Errorf("error: %v\n", err)
		}
	}
	zap.S().Infof("Broadcasted: %s to room %s\n", msg, r.Name)
}
