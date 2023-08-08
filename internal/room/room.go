package room

import (
	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
	"sync"
)

type Participant struct {
	conn *websocket.Conn
	ID   string
}

func newParticipant(c *websocket.Conn) *Participant {
	ID := c.Headers("X-Client-Id", "empty")
	return &Participant{conn: c, ID: ID}
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
	p := newParticipant(conn)
	r.participantMu.Lock()
	r.participants[p] = struct{}{}
	r.participantMu.Unlock()
	zap.S().Infof("Peer connected %s [%s] to room %s\n", conn.RemoteAddr(), p.ID, r.Name)
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
