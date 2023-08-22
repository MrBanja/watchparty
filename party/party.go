package party

import (
	"sync"

	"go.uber.org/zap"
)

type Party struct {
	partyMu sync.RWMutex
	party   map[string]*Room

	logger *zap.Logger
}

func New(logger *zap.Logger) *Party {
	return &Party{
		party:  make(map[string]*Room),
		logger: logger.Named("party"),
	}
}

func (p *Party) GetRoom(roomName string) *Room {
	p.partyMu.RLock()
	defer p.partyMu.RUnlock()
	return p.party[roomName]
}

func (p *Party) Create(roomName string) *Room {
	p.partyMu.Lock()
	defer p.partyMu.Unlock()
	room := NewRoom(roomName, p.logger)
	p.party[roomName] = room
	return room
}

func (p *Party) GetOrCreateRoom(roomName string) *Room {
	if room := p.GetRoom(roomName); room != nil {
		return room
	}
	return p.Create(roomName)
}
