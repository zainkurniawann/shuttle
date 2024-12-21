package services

import (
	// "strings"
	// "time"

	"database/sql"
	"log"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"github.com/jmoiron/sqlx"
)

type ChildernServiceInterface interface {
	GetAllChilderns(id string) ([]dto.StudentResponseDTO, int, error)
	GetSpecChildern(id string) (dto.StudentResponseDTO, error)
	UpdateChildern(tx *sqlx.Tx, id string, req dto.StudentRequestDTO, username string) error
}

type ChildernService struct {
	ChildernRepository repositories.ChildernRepositoryInterface
}

func NewChildernService(childernRepository repositories.ChildernRepositoryInterface) ChildernServiceInterface {
	return &ChildernService{
		ChildernRepository: childernRepository,
	}
}

func (service *ChildernService) GetAllChilderns(id string) ([]dto.StudentResponseDTO, int, error) {
    childerns, err := service.ChildernRepository.FetchAllChilderns(id)
    if err != nil {
        log.Println("Error fetching students from repository:", err)
        return nil, 0, err
    }

    var childernsDTO []dto.StudentResponseDTO

    for _, childern := range childerns {
        childernsDTO = append(childernsDTO, dto.StudentResponseDTO{
            UUID:          childern.UUID.String(),
            FirstName:     childern.FirstName,
            LastName:      childern.LastName,
            Grade:         childern.Grade,
            Gender:        childern.Gender,
            SchoolUUID:    childern.SchoolUUID.String(),
            ShuttleStatus: childern.ShuttleStatus.String, // Gunakan String
        })
    }

    total := len(childernsDTO)

    return childernsDTO, total, nil
}


func (service *ChildernService) GetSpecChildern(id string) (dto.StudentResponseDTO, error) {
	childern, err := service.ChildernRepository.FetchSpecChildern(id)
	if err != nil {
		return dto.StudentResponseDTO{}, err
	}

	studentDTO := dto.StudentResponseDTO{
		UUID:       childern.UUID.String(),
		FirstName:  childern.FirstName,
		LastName:   childern.LastName,
		Gender:     childern.Gender,
		Grade:      childern.Grade,
		SchoolUUID: childern.SchoolUUID.String(),
		CreatedAt:  safeTimeFormat(childern.CreatedAt),
		CreatedBy:  safeStringFormat(childern.CreatedBy),
		UpdatedAt:  safeTimeFormat(childern.UpdatedAt),
		UpdatedBy:  safeStringFormat(childern.UpdatedBy),
	}

	return studentDTO, nil
}

func (service *ChildernService) UpdateChildern(tx *sqlx.Tx, id string, req dto.StudentRequestDTO, username string) error {
	student := entity.Student{
		FirstName:  req.StudentFirstName,
		LastName:   req.StudentLastName,
		Gender:     string(req.StudentGender),
		UpdatedBy:  sql.NullString{String: username, Valid: username != ""},
	}

	err := service.ChildernRepository.UpdateChildern(tx, student, id)
	if err != nil {
		return err
	}

	return nil
}