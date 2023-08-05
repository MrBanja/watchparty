package service

import "github.com/gofiber/fiber/v2"

func (s *Service) GetMagnet(c *fiber.Ctx) error {
	return c.SendString("magnet:?xt=urn:xxxx")
}
