package service

import "github.com/gofiber/fiber/v2"

func (s *Service) Index(c *fiber.Ctx) error {
	return c.SendFile("./static/site/index.html")
}
