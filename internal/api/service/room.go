package service

import (
	fws "github.com/fasthttp/websocket"
	"github.com/gofiber/contrib/websocket"
	"go.uber.org/zap"
	"stream/internal/room"
)

func (s *Service) Room(c *websocket.Conn) {
	var (
		mt  int
		msg []byte
		err error
	)

	roomName := c.Params("id")
	s.mu.RLock()
	hub, ok := s.hubs[roomName]
	s.mu.RUnlock()
	if !ok {
		hub = room.New(roomName)
		s.mu.Lock()
		s.hubs[roomName] = hub
		s.mu.Unlock()
		zap.S().Infof("Room %s created", roomName)
	}

	participant := hub.AddParticipant(c)
	defer hub.RemoveParticipant(participant)

	for {
		if msg, err = participant.ReadMessage(); err != nil {
			if v, ok := err.(*fws.CloseError); ok {
				switch v.Code {
				case fws.CloseNormalClosure:
					zap.S().Infof("Peer disconnected %s from room %s normaly\n", c.RemoteAddr(), roomName)
				case fws.CloseGoingAway:
					zap.S().Infof("Peer disconnected %s to room %s gooing away\n", c.RemoteAddr(), roomName)
				default:
					zap.S().Warnf("Peer disconnected %s from room %s with error: %v\n", c.RemoteAddr(), roomName, err)
				}
			} else {
				zap.S().Errorf("Peer disconnected %s from room %s with unknown error: %v\n", c.RemoteAddr(), roomName, err)
			}
			break
		}
		zap.S().Infof("recv [%v]: %s: %s", mt, msg, c.RemoteAddr())

		hub.BroadcastExcept(string(msg), participant)
	}
}
