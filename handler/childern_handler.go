package handler

import (
	"log"
	"net/http"
	"strings"

	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ChildernHandlerInterface interface {
	GetAllChilderns(c *fiber.Ctx) error
	GetSpecChildern(c *fiber.Ctx) error
	UpdateChildern(c *fiber.Ctx) error
}

type ChildernHandler struct {
	DB              *sqlx.DB
	ChildernService services.ChildernServiceInterface
}

func NewChildernHandler(childernService services.ChildernServiceInterface) *ChildernHandler {
	return &ChildernHandler{
		ChildernService: childernService,
	}
}

func (handler *ChildernHandler) GetAllChilderns(c *fiber.Ctx) error {
	id, ok := c.Locals("userUUID").(string)
	if !ok || id == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "User UUID is missing or invalid",
		})
	}
	childernsDTO, total, err := handler.ChildernService.GetAllChilderns(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch students",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  childernsDTO,
		"total": total,
	})
}

func (handler *ChildernHandler) GetSpecChildern(c *fiber.Ctx) error {
	id := c.Params("id")

	log.Println("idddd", id)

	if _, err := uuid.Parse(id); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	studentDTO, err := handler.ChildernService.GetSpecChildern(id)
	if err != nil {
		log.Println("Error getting data from service:", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch data",
		})
	}

	log.Println("Student data:", studentDTO)

	return c.Status(http.StatusOK).JSON(studentDTO)
}

func (handler *ChildernHandler) UpdateChildern(c *fiber.Ctx) error {
	id := c.Params("id")
	username := c.Locals("user_name").(string)

	studentReqDTO := new(dto.StudentRequestDTO)
	if err := c.BodyParser(studentReqDTO); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, studentReqDTO); err != nil {
		return utils.BadRequestResponse(c, strings.ToUpper(err.Error()[0:1])+err.Error()[1:], nil)
	}

	tx, err := handler.DB.Beginx()
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to start transaction", nil)
	}
	defer tx.Rollback()

	existingStudent, err := handler.ChildernService.GetSpecChildern(id)
	if err != nil {
		logger.LogError(err, "Failed to fetch student", nil)
		return utils.NotFoundResponse(c, "Student not found", nil)
	}

	if existingStudent.UUID != id {
		return utils.NotFoundResponse(c, "Student UUID does not match", nil)
	}

	if err := handler.ChildernService.UpdateChildern(tx, id, *studentReqDTO, username); err != nil {
		logger.LogError(err, "Failed to update student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student updated successfully", nil)
}