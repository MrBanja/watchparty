package service

import "github.com/gofiber/fiber/v2"

func (s *Service) GetMagnet(c *fiber.Ctx) error {
	return c.SendString("magnet:?xt=urn:xxxx")
}

func (s *Service) GetStatusBadge(c *fiber.Ctx) error {
	roomName := c.Params("id")
	pID := c.Params("partID")

	s.mu.RLock()
	room, ok := s.hubs[roomName]
	s.mu.RUnlock()
	if !ok {
		return c.Render("status_badge", fiber.Map{"Address": s.address, "IsLost": true, "Text": "Empty room", "ID": roomName, "PartID": pID})
	}
	p := room.GetParticipantByID(pID)
	if p == nil {
		return c.Render("status_badge", fiber.Map{"Address": s.address, "IsLost": true, "Text": "Lost Connection", "ID": roomName})
	}
	return c.Render("status_badge", fiber.Map{"Address": s.address, "IsLost": false, "Text": "Connected", "ID": roomName, "PartID": p.ID})
}
