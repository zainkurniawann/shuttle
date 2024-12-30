package handler

import (
	"fmt"
	"log"
	"reflect"
	"shuttle/errors"
	"shuttle/logger"
	"shuttle/models/dto"
	"shuttle/services"
	"shuttle/utils"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type StudentHandlerInterface interface {
	GetAllStudentWithParents(c *fiber.Ctx) error
	GetSpecStudentWithParents(c *fiber.Ctx) error
	AddSchoolStudentWithParents(c *fiber.Ctx) error
	UpdateSchoolStudentWithParents(c *fiber.Ctx) error
	DeleteSchoolStudentWithParentsIfNeccessary(c *fiber.Ctx) error
}

type studentHandler struct {
	studentService services.StudentService
}

func NewStudentHttpHandler(studentService services.StudentService) StudentHandlerInterface {
	return &studentHandler{
		studentService: studentService,
	}
}

func (handler *studentHandler) GetAllStudentWithParents(c *fiber.Ctx) error {
	// Ambil schoolUUID dari token
	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		log.Println("ERROR: Token does not contain school UUID")
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Println("INFO: School UUID from token:", schoolUUIDStr)

	// Ambil halaman dan limit dari query params
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		log.Println("ERROR: Invalid page number:", c.Query("page"))
		return utils.BadRequestResponse(c, "Invalid page number", nil)
	}
	log.Println("INFO: Page number:", page)

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit < 1 {
		log.Println("ERROR: Invalid limit number:", c.Query("limit"))
		return utils.BadRequestResponse(c, "Invalid limit number", nil)
	}
	log.Println("INFO: Limit number:", limit)

	// Ambil sorting dari query params
	sortField := c.Query("sort_by", "student_id")
	sortDirection := c.Query("direction", "asc")

	if sortDirection != "asc" && sortDirection != "desc" {
		log.Println("ERROR: Invalid sort direction:", sortDirection)
		return utils.BadRequestResponse(c, "Invalid sort direction, use 'asc' or 'desc'", nil)
	}
	log.Println("INFO: Sort direction:", sortDirection)

	if !isValidSortFieldForStudents(sortField) {
		log.Println("ERROR: Invalid sort field:", sortField)
		return utils.BadRequestResponse(c, "Invalid sort field", nil)
	}
	log.Println("INFO: Sort field:", sortField)

	// Ambil data siswa dengan parents
	students, totalItems, err := handler.studentService.GetAllStudentsWithParents(page, limit, sortField, sortDirection, schoolUUIDStr)
	if err != nil {
		log.Println("ERROR: Failed to fetch paginated students:", err)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Println("INFO: Fetched students data successfully")

	// Hitung total pages
	totalPages := (totalItems + limit - 1) / limit
	log.Println("INFO: Total pages:", totalPages)

	// Validasi halaman yang diminta
	if page > totalPages {
		if totalItems > 0 {
			log.Println("ERROR: Page number out of range:", page)
			return utils.BadRequestResponse(c, "Page number out of range", nil)
		} else {
			page = 1
		}
	}
	log.Println("INFO: Adjusted page number:", page)

	start := (page-1)*limit + 1
	if totalItems == 0 || start > totalItems {
		start = 0
	}
	log.Println("INFO: Start index:", start)

	end := start + len(students) - 1
	if end > totalItems {
		end = totalItems
	}
	log.Println("INFO: End index:", end)

	if len(students) == 0 {
		start = 0
		end = 0
	}
	log.Println("INFO: Showing students:", start, "to", end)

	// Format response
	response := fiber.Map{
		"data": students,
		"meta": fiber.Map{
			"current_page":   page,
			"total_pages":    totalPages,
			"per_page_items": limit,
			"total_items":    totalItems,
			"showing":        fmt.Sprintf("Showing %d-%d of %d", start, end, totalItems),
		},
	}

	log.Println("INFO: Responding with student data")

	return utils.SuccessResponse(c, "Students fetched successfully", response)
}

func (handler *studentHandler) GetSpecStudentWithParents(c *fiber.Ctx) error {
	id := c.Params("id")

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	students, err := handler.studentService.GetSpecStudentWithParents(id, schoolUUIDStr)
	if err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to fetch students", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Students fetched successfully", students)
}

func (handler *studentHandler) AddSchoolStudentWithParents(c *fiber.Ctx) error {
	username, ok := c.Locals("user_name").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain username", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	student := new(dto.SchoolStudentParentRequestDTO)
	if err := c.BodyParser(student); err != nil {
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}

	if err := utils.ValidateStruct(c, student); err != nil {
		return utils.BadRequestResponse(c, err.Error(), nil)
	}

	// Validasi tambahan untuk field student_address dan student_pickup_point
	if student.Student.StudentAddress == "" {
		return utils.BadRequestResponse(c, "Address is required", nil)
	}

	// Validasi pickup point: pastikan ada latitude dan longitude
	if student.Student.StudentPickupPoint == nil || 
		student.Student.StudentPickupPoint["latitude"] == 0 || 
		student.Student.StudentPickupPoint["longitude"] == 0 {
		return utils.BadRequestResponse(c, "Valid latitude and longitude are required for pickup point", nil)
	}

	if reflect.DeepEqual(dto.UserRequestsDTO{}, student.Parent) {
		return utils.BadRequestResponse(c, "Parent details are required", nil)
	}

	if err := handler.studentService.AddSchoolStudentWithParents(*student, schoolUUIDStr, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to add student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student created successfully", nil)
}

func (handler *studentHandler) UpdateSchoolStudentWithParents(c *fiber.Ctx) error {
	username, ok := c.Locals("user_name").(string)
	if !ok {
		log.Println("ERROR: Token does not contain username")
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Println("INFO: Username extracted from token:", username)

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		log.Println("ERROR: Token does not contain school uuid")
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Println("INFO: School UUID extracted from token:", schoolUUIDStr)

	id := c.Params("id")
	log.Println("INFO: Received student ID from URL parameters:", id)

	student := new(dto.SchoolStudentParentRequestDTO)
	if err := c.BodyParser(student); err != nil {
		log.Println("ERROR: Failed to parse request body:", err)
		return utils.BadRequestResponse(c, "Invalid request data", nil)
	}
	log.Println("INFO: Request body successfully parsed:", student)

	if err := utils.ValidateStruct(c, student); err != nil {
		log.Println("ERROR: Validation failed for student data:", err)
		return utils.BadRequestResponse(c, err.Error(), nil)
	}
	log.Println("INFO: Student data successfully validated")

	// Validasi tambahan untuk field student_address dan student_pickup_point
	if student.Student.StudentAddress == "" {
		log.Println("ERROR: Student address is missing")
		return utils.BadRequestResponse(c, "Address is required", nil)
	}
	log.Println("INFO: Student address is provided:", student.Student.StudentAddress)

	// Validasi pickup point: pastikan ada latitude dan longitude
	if student.Student.StudentPickupPoint == nil || 
		student.Student.StudentPickupPoint["latitude"] == 0 || 
		student.Student.StudentPickupPoint["longitude"] == 0 {
		log.Println("ERROR: Invalid pickup point (latitude or longitude missing)")
		return utils.BadRequestResponse(c, "Valid latitude and longitude are required for pickup point", nil)
	}
	log.Println("INFO: Pickup point validated:", student.Student.StudentPickupPoint)

	if reflect.DeepEqual(dto.UserRequestsDTO{}, student.Parent) {
		log.Println("ERROR: Parent details are missing")
		return utils.BadRequestResponse(c, "Parent details are required", nil)
	}
	log.Println("INFO: Parent details provided:", student.Parent)

	// Call to service layer to update student data
	log.Println("INFO: Updating student data in service layer:", map[string]interface{}{"studentID": id, "schoolUUID": schoolUUIDStr, "username": username})
	if err := handler.studentService.UpdateSchoolStudentWithParents(id, *student, schoolUUIDStr, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			log.Println("ERROR: Custom error occurred while updating student:", err)
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		log.Println("ERROR: Failed to update student in service layer:", err)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}
	log.Println("INFO: Student data successfully updated:", id)

	return utils.SuccessResponse(c, "Student updated successfully", nil)
}



func (handler *studentHandler) DeleteSchoolStudentWithParentsIfNeccessary(c *fiber.Ctx) error {
	id := c.Params("id")

	schoolUUIDStr, ok := c.Locals("schoolUUID").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain school uuid", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	username, ok := c.Locals("user_name").(string)
	if !ok {
		logger.LogError(nil, "Token does not contain username", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	if err := handler.studentService.DeleteSchoolStudentWithParentsIfNeccessary(id, schoolUUIDStr, username); err != nil {
		if customErr, ok := err.(*errors.CustomError); ok {
			return utils.ErrorResponse(c, customErr.StatusCode, strings.ToUpper(string(customErr.Message[0]))+customErr.Message[1:], nil)
		}
		logger.LogError(err, "Failed to delete student", nil)
		return utils.InternalServerErrorResponse(c, "Something went wrong, please try again later", nil)
	}

	return utils.SuccessResponse(c, "Student deleted successfully", nil)
}

func isValidSortFieldForStudents(field string) bool {
	allowedFields := map[string]bool{
		"student_id":        true,
		"student_grade":     true,
		"student_first_name": true,
		"student_last_name":  true,
	}
	return allowedFields[field]
}