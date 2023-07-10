package service

import (
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
	hub, ok := s.hubs[roomName]
	if !ok {
		hub = room.New(roomName)
		s.hubs[roomName] = hub
		zap.S().Infof("Room %s created", roomName)
	}

	participant := hub.AddParticipant(c)
	defer hub.RemoveParticipant(participant)

	for {
		if msg, err = participant.ReadMessage(); err != nil {
			zap.S().Error("read:", err)
			break
		}
		zap.S().Infof("recv [%v]: %s", mt, msg)

		hub.Broadcast(string(msg))
	}
}
