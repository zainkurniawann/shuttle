package handler

import (
	"net/http"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strings"
	"errors"
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ShuttleHandler struct {
	ShuttleService services.ShuttleServiceInterface
	DB             *sqlx.DB
}

func NewShuttleHandler(shuttleService services.ShuttleServiceInterface) *ShuttleHandler {
	return &ShuttleHandler{
		ShuttleService: shuttleService,
	}
}

func (h *ShuttleHandler) GetShuttleTrackByParent(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
	}

	parentUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
	}

	shuttles, err := h.ShuttleService.GetShuttleTrackByParent(parentUUID)
	if err != nil {
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}

	return c.Status(http.StatusOK).JSON(shuttles)
}

func (h *ShuttleHandler) GetSpecShuttle(c *fiber.Ctx) error {
	shuttleUUIDParam := c.Params("id")
	if shuttleUUIDParam == "" {
		return utils.BadRequestResponse(c, "Missing shuttle UUID in request URL", nil)
	}

	shuttleUUID, err := uuid.Parse(shuttleUUIDParam)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid shuttle UUID format", nil)
	}

	shuttle, err := h.ShuttleService.GetSpecShuttle(shuttleUUID)
	if err != nil {
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}

	if len(shuttle) == 0 {
		return utils.NotFoundResponse(c, "Shuttle data not found", nil)
	}

	return c.Status(http.StatusOK).JSON(shuttle)
}

func (h *ShuttleHandler) AddShuttle(c *fiber.Ctx) error {
	userUUID, ok := c.Locals("userUUID").(string)
	if !ok || userUUID == "" {
		return utils.BadRequestResponse(c, "Invalid or missing userUUID", nil)
	}
	username := c.Locals("user_name").(string)
	driverUUID, err := uuid.Parse(userUUID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid userUUID format", nil)
	}

	shuttleReq := new(dto.ShuttleRequest)
	if err := c.BodyParser(shuttleReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if shuttleReq.Status == "" {
		shuttleReq.Status = "waiting"
	}

	if err := utils.ValidateStruct(c, shuttleReq); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	if err := h.ShuttleService.AddShuttle(*shuttleReq, driverUUID.String(), username); err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to add shuttle", nil)
	}

	return utils.SuccessResponse(c, "Shuttle added successfully", nil)
}

func (h *ShuttleHandler) EditShuttle(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return utils.BadRequestResponse(c, "Missing shuttleUUID in URL", nil)
	}

	var statusReq struct {
		Status string `json:"status" validate:"required"`
	}
	if err := c.BodyParser(&statusReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", nil)
	}

	if err := utils.ValidateStruct(c, statusReq); err != nil {
		return utils.BadRequestResponse(c, "Invalid status: "+err.Error(), nil)
	}

	if err := h.ShuttleService.EditShuttleStatus(id, statusReq.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return utils.NotFoundResponse(c, "Shuttle not found", nil)
		}
		return utils.InternalServerErrorResponse(c, "Failed to edit shuttle", nil)
	}

	return utils.SuccessResponse(c, "Shuttle status updated successfully", nil)
}
