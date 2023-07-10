package service

import "stream/internal/room"

type Service struct {
	hubs map[string]*room.Room
}

func New() *Service {
	return &Service{
		hubs: make(map[string]*room.Room),
	}
}
