package handler

import (
	"shuttle/services"
	"github.com/gofiber/fiber/v2"
)

type ShuttleNotificationHandler struct {
	ShuttleNotificationService *services.ShuttleNotificationService
}

func NewShuttleNotificationHandler(service *services.ShuttleNotificationService) *ShuttleNotificationHandler {
	return &ShuttleNotificationHandler{ShuttleNotificationService: service}
}

func (h *ShuttleNotificationHandler) StartShuttleMonitoring(c *fiber.Ctx) error {
	go h.ShuttleNotificationService.MonitorShuttleStatus()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shuttle status monitoring started",
	})
}
