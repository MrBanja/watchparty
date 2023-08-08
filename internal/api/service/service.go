package service

import (
	"stream/internal/room"
	"sync"
)

type Service struct {
	address string

	mu   sync.RWMutex
	hubs map[string]*room.Room
}

func New(address string) *Service {
	return &Service{
		address: address,
		hubs:    make(map[string]*room.Room),
	}
}
