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
	GetAllRoutes(c *fiber.Ctx) error
	GetSpecRoute(c *fiber.Ctx) error
	GetAllRoutesByDriver(c *fiber.Ctx) error
	GetSpecRouteByDriver(c *fiber.Ctx) error 
	AddRoute(c *fiber.Ctx) error
	UpdateRoute(c *fiber.Ctx) error
	DeleteRoute(c *fiber.Ctx) error
}

type routeHandler struct {
	routeService services.RouteServiceInterface
}

func NewRouteHttpHandler(routeService services.RouteServiceInterface) RouteHandlerInterface {
	return &routeHandler{
		routeService: routeService,
	}
}

func (handler *routeHandler) GetAllRoutes(c *fiber.Ctx) error {
	log.Println("Starting GetAllRoutes handler")

	// Memanggil service untuk mendapatkan semua routes
	routes, err := handler.routeService.GetAllRoutes()
	if err != nil {
		log.Printf("Failed to get routes: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch routes",
		})
	}

	// Mengembalikan response dengan data routes
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"routes": routes,
	})
}

func (handler *routeHandler) GetSpecRoute(c *fiber.Ctx) error {
	log.Println("Starting GetSpec handler")

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

	// Mengambil username dari context (misalnya dari token)
	username, ok := c.Locals("user_name").(string)
	if !ok {
		log.Println("Token does not contain username")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token does not contain username",
		})
	}
	log.Printf("Username: %s", username)

	// Panggil service untuk mendapatkan detail route
	route, err := handler.routeService.GetSpecRoute(routeUUID)
	if err != nil {
		log.Printf("Failed to get route by UUID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"route": route,
	})
}

func (handler *routeHandler) GetAllRoutesByDriver(c *fiber.Ctx) error {
	// Ambil driverUUID dari context (dari token)
	driverUUID, ok := c.Locals("userUUID").(string)
	if !ok {
		log.Println("Token does not contain driver UUID")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token does not contain driver UUID",
		})
	}

	// Validasi format UUID
	_, err := uuid.Parse(driverUUID)
	if err != nil {
		log.Printf("Invalid UUID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UUID format",
		})
	}

	// Panggil service untuk mendapatkan routes berdasarkan driver UUID
	routes, err := handler.routeService.GetAllRoutesByDriver(driverUUID)
	if err != nil {
		log.Printf("Failed to get routes by driver UUID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch routes",
		})
	}

	// Mengembalikan response dengan data routes
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"routes": routes,
	})
}

func (handler *routeHandler) GetSpecRouteByDriver(c *fiber.Ctx) error {
	// Ambil driverUUID dari token
	driverUUID, ok := c.Locals("userUUID").(string)
	if !ok || driverUUID == "" {
		log.Println("Token does not contain driver UUID")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// Ambil studentUUID dari URL parameter
	studentUUID := c.Params("id")
	if studentUUID == "" {
		log.Println("Student UUID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Student UUID is required",
		})
	}

	// Panggil service untuk mendapatkan data spesifik route
	route, err := handler.routeService.GetSpecRouteByDriver(driverUUID, studentUUID)
	if err != nil {
		log.Printf("Failed to get specific route: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch specific route",
		})
	}

	// Mengembalikan response dengan data spesifik route
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"route": route,
	})
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

func (handler *routeHandler) DeleteRoute(c *fiber.Ctx) error {
	log.Println("Starting DeleteRoute handler")

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

	// Mengambil username dari context (misalnya dari token)
	username, ok := c.Locals("user_name").(string)
	if !ok {
		log.Println("Token does not contain username")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Token does not contain username",
		})
	}
	log.Printf("Username: %s", username)

	// Panggil service untuk delete route
	if err := handler.routeService.DeleteRoute(routeUUID, username); err != nil {
		log.Printf("Failed to delete route: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Return success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Route deleted successfully",
	})
}
