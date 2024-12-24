package handler

import (
	"log"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type RouteHandlerInterface interface {
	AddRoute(c *fiber.Ctx) error
	UpdateRoute(c *fiber.Ctx) error
}

type routeHandler struct {
	routeService services.RouteServiceInterface
}

func NewRouteHttpHandler(routeService services.RouteServiceInterface) RouteHandlerInterface {
	return &routeHandler{
		routeService: routeService,
	}
}

func (handler *routeHandler) AddRoute(c *fiber.Ctx) error {
	log.Println("Starting AddRoute handler")

	// Ambil schoolUUID dari token
	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		log.Println("Token does not contain school UUID")
		return utils.InternalServerErrorResponse(c, "Token does not contain school UUID", nil)
	}
	log.Printf("School UUID: %s", schoolUUIDStr)

	// Ambil username dari token
	username, ok := c.Locals("user_name").(string)
	if !ok {
		log.Println("Token does not contain username")
		return utils.InternalServerErrorResponse(c, "Token does not contain username", nil)
	}
	log.Printf("Username: %s", username)

	// Parsing body ke DTO
	route := new(dto.RouteRequestDTO)
	if err := c.BodyParser(route); err != nil {
		log.Printf("Error parsing body: %v", err)
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}
	log.Printf("Route DTO parsed: %+v", route)

	// Validasi DTO
	if err := utils.ValidateStruct(c, route); err != nil {
		log.Printf("Validation failed: %v", err)
		return utils.BadRequestResponse(c, err.Error(), nil)
	}
	log.Println("Validation passed")

	// Panggil service untuk menambahkan route
	if err := handler.routeService.AddRoute(*route, schoolUUIDStr, username); err != nil {
		log.Printf("Error adding route: %v", err)
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}
	log.Println("Route added successfully")

	return utils.SuccessResponse(c, "Route added successfully", nil)
}

func (handler *routeHandler) UpdateRoute(c *fiber.Ctx) error {
	log.Println("Starting UpdateRoute handler")

	// Ambil route_uuid dari URL
	routeUUID := c.Params("id")
	log.Printf("Extracted route_uuid from URL: %s", routeUUID)

	// Validasi UUID
	if _, err := uuid.Parse(routeUUID); err != nil {
		log.Printf("Invalid route UUID: %s", routeUUID)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid route UUID",
		})
	}
	log.Println("Route UUID validation passed")

	// Ambil data dari request body
	var updateDTO dto.RouteRequestDTO
	if err := c.BodyParser(&updateDTO); err != nil {
		log.Println("Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}
	log.Printf("Parsed request body: %+v", updateDTO)

	// Ambil username dari context (misalnya dari JWT)
	username, ok := c.Locals("user_name").(string)
	if !ok || username == "" {
		log.Println("Failed to retrieve username from context")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	log.Printf("Retrieved username from context: %s", username)

	// Panggil service untuk update
	log.Println("Calling routeService.UpdateRoute")
	if err := handler.routeService.UpdateRoute(routeUUID, updateDTO, username); err != nil {
		log.Printf("Failed to update route: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Println("Route updated successfully")

	// Response success
	log.Println("Returning success response")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Route updated successfully",
	})
}
