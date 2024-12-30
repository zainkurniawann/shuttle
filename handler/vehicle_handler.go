package handler

import (
	"fmt"
	"log"
	"net/http"
	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type VehicleHandlerInterface interface {
	GetAllVehicles(c *fiber.Ctx) error
	GetSpecVehicle(c *fiber.Ctx) error
	AddVehicle(c *fiber.Ctx) error
	UpdateVehicle(c *fiber.Ctx) error
	DeleteVehicle(c *fiber.Ctx) error
}

type vehicleHandler struct {
	vehicleService services.VehicleService
}

func NewVehicleHttpHandler(vehicleService services.VehicleService) VehicleHandlerInterface {
	return &vehicleHandler{
		vehicleService: vehicleService,
	}
}

func (handler *vehicleHandler) GetAllVehicles(c *fiber.Ctx) error {
    role := c.Locals("role_code").(string)
    
    // Ambil schoolUUID dari context
    schoolUUID, ok := c.Locals("schoolUUID").(string)
    if !ok {
        return utils.BadRequestResponse(c, "Invalid token or schoolUUID", nil)
    }

    // Ambil page, limit, sortField, dan sortDirection dari query
    page, err := strconv.Atoi(c.Query("page", "1"))
    if err != nil || page < 1 {
        return utils.BadRequestResponse(c, "Invalid page number", nil)
    }

    limit, err := strconv.Atoi(c.Query("limit", "10"))
    if err != nil || limit < 1 {
        return utils.BadRequestResponse(c, "Invalid limit number", nil)
    }

    sortField := c.Query("sort_by", "vehicle_id")
    sortDirection := c.Query("direction", "asc")

    if sortDirection != "asc" && sortDirection != "desc" {
        return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
    }

    if !isValidSortFieldForVehicles(sortField) {
        return utils.BadRequestResponse(c, "Invalid sort field", nil)
    }

    var vehicles []dto.VehicleResponseDTO
    var totalItems int

    switch role {
    case "SA": // Super Admin
        vehicles, totalItems, err = handler.vehicleService.GetAllVehicles(page, limit, sortField, sortDirection)
        if err != nil {
            logger.LogError(err, "Failed to fetch vehicles for Super Admin", nil)
            return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
        }
    case "AS": // Admin Sekolah
        vehicles, totalItems, err = handler.vehicleService.GetAllVehiclesForSchool(page, limit, sortField, sortDirection, schoolUUID)
        if err != nil {
            logger.LogError(err, "Failed to fetch vehicles for Admin School", nil)
            return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
        }
    default:
        return utils.BadRequestResponse(c, "Invalid role", nil)
    }

    // Hitung total halaman
    totalPages := (totalItems + limit - 1) / limit
    if page > totalPages {
        if totalItems > 0 {
            return utils.BadRequestResponse(c, "Page number out of range", nil)
        } else {
            page = 1
        }
    }

    start := (page-1)*limit + 1
    if totalItems == 0 || start > totalItems {
        start = 0
    }

    end := start + len(vehicles) - 1
    if end > totalItems {
        end = totalItems
    }

    if len(vehicles) == 0 {
        start = 0
        end = 0
    }

    response := fiber.Map{
        "data": vehicles,
        "meta": fiber.Map{
            "current_page":   page,
            "total_pages":    totalPages,
            "per_page_items": limit,
            "total_items":    totalItems,
            "showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
        },
    }

    return utils.SuccessResponse(c, "Vehicles fetched successfully", response)
}

func (handler *vehicleHandler) GetSpecVehicle(c *fiber.Ctx) error {
	id := c.Params("id")
	vehicle, err := handler.vehicleService.GetSpecVehicle(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch specific vehicle", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Vehicle fetched successfully", vehicle)
}

func (handler *vehicleHandler) AddVehicle(c *fiber.Ctx) error {
    log.Println("Start processing AddVehicle request")

    // Parsing body request ke DTO
    vehicle := new(dto.VehicleRequestDTO)
    if err := c.BodyParser(vehicle); err != nil {
        log.Println("Error parsing request body:", err)
        return utils.BadRequestResponse(c, "Invalid request data", nil)
    }
    log.Println("Request body parsed successfully:", vehicle)

    // Validasi request
    if err := utils.ValidateStruct(c, vehicle); err != nil {
        log.Println("Validation error:", err)
        return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
    }
    log.Println("Request body validation passed")

    // Ambil role dan user_id dari token
    role, ok := c.Locals("role_code").(string)
    if !ok || role == "" {
        log.Println("User role missing or invalid")
        return utils.UnauthorizedResponse(c, "Invalid user role", nil)
    }
    log.Println("User role retrieved from token:", role)

    userID, ok := c.Locals("userUUID").(string)
    if !ok || userID == "" {
        log.Println("User ID missing in token")
        return utils.UnauthorizedResponse(c, "Invalid user ID", nil)
    }
    log.Println("User ID retrieved from token:", userID)

    // Ambil schoolUUID dari context
    schoolUUID, ok := c.Locals("schoolUUID").(string)
    if !ok || schoolUUID == "" {
        log.Println("School UUID missing or invalid in context")
        return utils.UnauthorizedResponse(c, "School UUID is missing or invalid", nil)
    }
    log.Println("School UUID retrieved from context:", schoolUUID)

    // Panggil service untuk menambahkan vehicle
    log.Println("Calling AddVehicle service")
    if err := handler.vehicleService.AddVehicle(*vehicle, role, schoolUUID); err != nil {
        if customErr, ok := err.(*errors.CustomError); ok {
            log.Println("Error from AddVehicle service:", customErr.Message)
            return utils.ErrorResponse(c, customErr.StatusCode, customErr.Message, nil)
        }
        log.Println("Unexpected error from AddVehicle service:", err)
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong, please try again later", nil)
    }
    log.Println("Vehicle successfully added")

    // Berhasil
    return utils.SuccessResponse(c, "Vehicle created successfully", nil)
}

func (handler *vehicleHandler) UpdateVehicle(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	vehicle := new(dto.VehicleRequestDTO)
	if err := c.BodyParser(vehicle); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, vehicle); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	if err := handler.vehicleService.UpdateVehicle(id, *vehicle, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Vehicle updated successfully", nil)
}

func (handler *vehicleHandler) DeleteVehicle(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	if err := handler.vehicleService.DeleteVehicle(id, username); err != nil {
		logger.LogError(err, "Failed to delete vehicle", nil)
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Vehicle deleted successfully", nil)
}

func isValidSortFieldForVehicles(field string) bool {
	allowedFields := map[string]bool{
		"vehicle_id":     true,
		"vehicle_name":   true,
		"vehicle_number": true,
		"vehicle_type":   true,
		"vehicle_color":  true,
		"vehicle_seats":  true,
		"vehicle_status": true,
	}

	return allowedFields[field]
}
