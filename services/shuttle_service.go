package services

import (
	"database/sql"
	"fmt"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type ShuttleServiceInterface interface {
	GetShuttleTrackByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error)
	GetAllShuttleByParent(parentUUID uuid.UUID) ([]dto.ShuttleAllResponse, error)
	GetSpecShuttle(shuttleUUID uuid.UUID) ([]dto.ShuttleSpecResponse, error)
	AddShuttle(req dto.ShuttleRequest, driverUUID, createdBy string) error
	EditShuttleStatus(shuttleUUID, status string) error
}

type ShuttleService struct {
	shuttleRepository repositories.ShuttleRepositoryInterface
}

func NewShuttleService(shuttleRepository repositories.ShuttleRepositoryInterface) ShuttleServiceInterface {
	return &ShuttleService{
		shuttleRepository: shuttleRepository,
	}
}

func (s *ShuttleService) GetShuttleTrackByParent(parentUUID uuid.UUID) ([]dto.ShuttleResponse, error) {
	shuttles, err := s.shuttleRepository.FetchShuttleTrackByParent(parentUUID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ShuttleResponse, 0, len(shuttles))
	for _, shuttle := range shuttles {
		response := &dto.ShuttleResponse{
			StudentUUID:    shuttle.StudentUUID,
			ShuttleUUID:    shuttle.ShuttleUUID,
			StudentName:    shuttle.StudentName,
			StudentLastName: shuttle.StudentLastName,
			ParentUUID:     shuttle.ParentUUID,
			SchoolUUID:     shuttle.SchoolUUID,
			SchoolName:     shuttle.SchoolName,
			ShuttleStatus:  shuttle.ShuttleStatus,
			CreatedAt:      shuttle.CreatedAt,
			CurrentDate:    shuttle.CurrentDate,
		}
		responses = append(responses, *response)
	}

	return responses, nil
}

func (s *ShuttleService) GetAllShuttleByParent(parentUUID uuid.UUID) ([]dto.ShuttleAllResponse, error) {
	// Fetch data from the repository
	shuttles, err := s.shuttleRepository.FetchAllShuttleByParent(parentUUID)
	if err != nil {
		return nil, err
	}

	// Transform the data if needed (DTO is already in the required format)
	responses := make([]dto.ShuttleAllResponse, len(shuttles))
	for i, shuttle := range shuttles {
		responses[i] = dto.ShuttleAllResponse{
			StudentUUID:     shuttle.StudentUUID,
			Status:          shuttle.Status,
			StudentFirstName: shuttle.StudentFirstName,
			StudentLastName:  shuttle.StudentLastName,
			StudentGrade:     shuttle.StudentGrade,
			StudentGender:    shuttle.StudentGender,
			ParentUUID:       shuttle.ParentUUID,
			SchoolUUID:       shuttle.SchoolUUID,
			SchoolName:       shuttle.SchoolName,
			CreatedAt:        shuttle.CreatedAt,
			UpdatedAt:        shuttle.UpdatedAt,
		}
	}

	return responses, nil
}


func (s *ShuttleService) GetSpecShuttle(shuttleUUID uuid.UUID) ([]dto.ShuttleSpecResponse, error) {
	shuttles, err := s.shuttleRepository.GetSpecShuttle(shuttleUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shuttle data: %w", err)
	}

	responses := make([]dto.ShuttleSpecResponse, 0, len(shuttles))
	for _, shuttle := range shuttles {
		response := dto.ShuttleSpecResponse{
			StudentUUID:       shuttle.StudentUUID,
			ShuttleUUID:       shuttle.ShuttleUUID,
			StudentFirstName:  shuttle.StudentFirstName,
			StudentLastName:   shuttle.StudentLastName,
			StudentPickupPoint: shuttle.StudentPickupPoint,
			ParentUUID:        shuttle.ParentUUID,
			SchoolUUID:        shuttle.SchoolUUID,
			SchoolName:        shuttle.SchoolName,
			SchoolPoint:       shuttle.SchoolPoint,
			ShuttleStatus:     shuttle.ShuttleStatus,
			CreatedAt:         shuttle.CreatedAt,
			CurrentDate:       shuttle.CurrentDate,
			DriverUUID:        shuttle.DriverUUID,
			DriverUsername:    shuttle.DriverUsername,
			DriverFirstName:   shuttle.DriverFirstName,
			DriverLastName:    shuttle.DriverLastName,
			DriverGender:      shuttle.DriverGender,
			VehicleUUID:       shuttle.VehicleUUID,
			VehicleName:       shuttle.VehicleName,
			VehicleType:       shuttle.VehicleType,
			VehicleColor:      shuttle.VehicleColor,
			VehicleNumber:     shuttle.VehicleNumber,
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *ShuttleService) AddShuttle(req dto.ShuttleRequest, driverUUID, createdBy string) error {
	studentUUID, err := uuid.Parse(req.StudentUUID)
	if err != nil {
		return err
	}

	driverUUIDParsed, err := uuid.Parse(driverUUID)
	if err != nil {
		return err
	}

	if req.Status == "" {
		req.Status = "waiting"
	}

	shuttle := entity.Shuttle{
		ShuttleID:   time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		ShuttleUUID: uuid.New(),
		StudentUUID: studentUUID,
		DriverUUID:  driverUUIDParsed,
		Status:      req.Status,
		CreatedAt:   sql.NullTime{Time: time.Now(), Valid: true},
	}

	err = s.shuttleRepository.SaveShuttle(shuttle)
	if err != nil {
		return err
	}

	return nil
}

func (s *ShuttleService) EditShuttleStatus(shuttleUUID, status string) error {
	shuttleUUIDParsed, err := uuid.Parse(shuttleUUID)
	if err != nil {
		return err
	}

	if err := s.shuttleRepository.UpdateShuttleStatus(shuttleUUIDParsed, status); err != nil {
		return err
	}

	return nil
}
