package services

import (
	"log"
	"shuttle/errors"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type VehicleServiceInterface interface {
	GetSpecVehicle(uuid string) (dto.VehicleResponseDTO, error)
	GetAllVehicles(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error)
	AddVehicle(req dto.VehicleRequestDTO) error
	GetSchoolUUIDByUserID(userID string) (string, error)
	UpdateVehicle(id string, req dto.VehicleRequestDTO, username string) error
	DeleteVehicle(id string, username string) error
}

type VehicleService struct {
	vehicleRepository repositories.VehicleRepositoryInterface
}

func NewVehicleService(vehicleRepository repositories.VehicleRepositoryInterface) VehicleService {
	return VehicleService{
		vehicleRepository: vehicleRepository,
	}
}

func (service *VehicleService) GetAllVehicles(page, limit int, sortField, sortDirection string) ([]dto.VehicleResponseDTO, int, error) {
	offset := (page - 1) * limit

	vehicles, school, driver, err := service.vehicleRepository.FetchAllVehicles(offset, limit, sortField, sortDirection)
	if err != nil {
		return nil, 0, err
	}

	total, err := service.vehicleRepository.CountVehicles()
	if err != nil {
		return nil, 0, err
	}

	var vehiclesDTO []dto.VehicleResponseDTO
	for _, vehicle := range vehicles {

		var schoolName string
		if vehicle.SchoolUUID == nil || school[vehicle.SchoolUUID.String()].UUID == uuid.Nil {
			schoolName = "N/A"
		} else {
			schoolName = school[vehicle.SchoolUUID.String()].Name
		}

		var driverName string
		if vehicle.DriverUUID == nil || driver[vehicle.DriverUUID.String()].UserUUID == uuid.Nil {
			driverName = "N/A"
		} else {
			driverName = driver[vehicle.DriverUUID.String()].FirstName + " " + driver[vehicle.DriverUUID.String()].LastName
		}

		vehiclesDTO = append(vehiclesDTO, dto.VehicleResponseDTO{
			UUID:       vehicle.UUID.String(),
			SchoolName: schoolName,
			DriverName: driverName,
			Name:       vehicle.VehicleName,
			Number:     vehicle.VehicleNumber,
			Type:       vehicle.VehicleType,
			Color:      vehicle.VehicleColor,
			Seats:      vehicle.VehicleSeats,
			Status:     vehicle.VehicleStatus,
			CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
		})
	}

	return vehiclesDTO, total, nil
}

func (service *VehicleService) GetAllVehiclesForSchool(page, limit int, sortField, sortDirection, schoolUUID string) ([]dto.VehicleResponseDTO, int, error) {
    offset := (page - 1) * limit

    // Modifikasi query untuk memasukkan schoolUUID
    vehicles, school, driver, err := service.vehicleRepository.FetchAllVehiclesForSchool(offset, limit, sortField, sortDirection, schoolUUID)
    if err != nil {
        return nil, 0, err
    }

	total, err := service.vehicleRepository.CountVehiclesForSchool(schoolUUID)
    if err != nil {
        return nil, 0, err
    }

    var vehiclesDTO []dto.VehicleResponseDTO
    for _, vehicle := range vehicles {

        var schoolName string
        if vehicle.SchoolUUID == nil || school[vehicle.SchoolUUID.String()].UUID == uuid.Nil {
            schoolName = "N/A"
        } else {
            schoolName = school[vehicle.SchoolUUID.String()].Name
        }

        var driverName string
        if vehicle.DriverUUID == nil || driver[vehicle.DriverUUID.String()].UserUUID == uuid.Nil {
            driverName = "N/A"
        } else {
            driverName = driver[vehicle.DriverUUID.String()].FirstName + " " + driver[vehicle.DriverUUID.String()].LastName
        }

        vehiclesDTO = append(vehiclesDTO, dto.VehicleResponseDTO{
            UUID:       vehicle.UUID.String(),
            SchoolName: schoolName,
            DriverName: driverName,
            Name:       vehicle.VehicleName,
            Number:     vehicle.VehicleNumber,
            Type:       vehicle.VehicleType,
            Color:      vehicle.VehicleColor,
            Seats:      vehicle.VehicleSeats,
            Status:     vehicle.VehicleStatus,
            CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
        })
    }

    return vehiclesDTO, total, nil
}

func (service *VehicleService) GetSpecVehicle(id string) (dto.VehicleResponseDTO, error) {
	vehicle, school, driver, err := service.vehicleRepository.FetchSpecVehicle(id)
	if err != nil {
		return dto.VehicleResponseDTO{}, err
	}

	var schoolUUID, schoolName string
	if vehicle.SchoolUUID == nil {
		schoolUUID = "N/A"
		schoolName = "N/A"
	} else if vehicle.SchoolUUID != nil {
		schoolUUID = vehicle.SchoolUUID.String()
		schoolName = school.Name
	}

	var driverUUID, driverName string
	if driver.UserUUID == uuid.Nil {
		driverUUID = "N/A"
		driverName = "N/A"
	} else if driver.UserUUID != uuid.Nil {
		driverUUID = vehicle.DriverUUID.String()
		driverName = driver.FirstName + " " + driver.LastName
	}

	vehicleDTO := dto.VehicleResponseDTO{
		UUID:       vehicle.UUID.String(),
		SchoolUUID: schoolUUID,
		SchoolName: schoolName,
		DriverUUID: driverUUID,
		DriverName: driverName,
		Name:       vehicle.VehicleName,
		Number:     vehicle.VehicleNumber,
		Type:       vehicle.VehicleType,
		Color:      vehicle.VehicleColor,
		Seats:      vehicle.VehicleSeats,
		Status:     vehicle.VehicleStatus,
		CreatedAt:  safeTimeFormat(vehicle.CreatedAt),
		CreatedBy:  safeStringFormat(vehicle.CreatedBy),
		UpdatedAt:  safeTimeFormat(vehicle.UpdatedAt),
		UpdatedBy:  safeStringFormat(vehicle.UpdatedBy),
	}

	return vehicleDTO, nil
}

func (service *VehicleService) GetSchoolUUIDByUserID(userID string) (string, error) {
	log.Println("Calling repository to fetch school_uuid for user_id:", userID)
	schoolUUID, err := service.vehicleRepository.GetSchoolUUIDByUserID(userID)
	if err != nil {
		log.Println("Error from repository while fetching school_uuid:", err)
		return "", err
	}
	log.Println("School UUID successfully retrieved from repository:", schoolUUID)
	return schoolUUID, nil
}


func (service *VehicleService) AddVehicle(req dto.VehicleRequestDTO, role, schoolUUID string) error {
    log.Println("Start adding vehicle")

    // Membuat entitas vehicle
    vehicle := entity.Vehicle{
        ID:            time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
        UUID:          uuid.New(),
        VehicleName:   req.Name,
        VehicleNumber: req.Number,
        VehicleType:   req.Type,
        VehicleColor:  req.Color,
        VehicleSeats:  req.Seats,
        VehicleStatus: req.Status,
    }
    log.Printf("Vehicle entity created: %+v\n", vehicle)

    // Gunakan schoolUUID yang sudah ada di context
    if role == "AS" {
        log.Println("Role is schooladmin, using school_uuid from token")
        if schoolUUID != "" {
            schoolUUIDParsed, err := uuid.Parse(schoolUUID)
            if err != nil {
                log.Println("Error parsing school UUID:", err)
                return errors.New("Invalid school UUID", 400)
            }
            vehicle.SchoolUUID = &schoolUUIDParsed
            log.Println("School UUID parsed and assigned to vehicle:", schoolUUIDParsed)
        } else {
            log.Println("School UUID is required for schooladmin")
            return errors.New("School UUID is required for schooladmin", 400)
        }
    } else {
        log.Println("Non-schooladmin role, using provided school UUID or nil")
    }

    // Cek apakah vehicle number sudah ada
    log.Println("Checking if vehicle number exists:", vehicle.VehicleNumber)
    isExistingVehicleNumber, err := service.vehicleRepository.CheckVehicleNumberExists("", vehicle.VehicleNumber)
    if err != nil {
        log.Println("Error checking vehicle number:", err)
        return err
    }

    if isExistingVehicleNumber {
        log.Println("Vehicle number already exists:", vehicle.VehicleNumber)
        return errors.New("Vehicle number already exists", 400)
    }
    log.Println("Vehicle number is unique, proceeding to save")

    // Simpan kendaraan
    log.Println("Saving vehicle to database")
    err = service.vehicleRepository.SaveVehicle(vehicle)
    if err != nil {
        log.Println("Error saving vehicle:", err)
        return err
    }
    log.Println("Vehicle saved successfully")
    return nil
}

func (service *VehicleService) UpdateVehicle(id string, req dto.VehicleRequestDTO, username string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	vehicle := entity.Vehicle{
		UUID:          parsedUUID,
		VehicleName:   req.Name,
		VehicleNumber: req.Number,
		VehicleType:   req.Type,
		VehicleColor:  req.Color,
		VehicleSeats:  req.Seats,
		VehicleStatus: req.Status,
		UpdatedAt:     toNullTime(time.Now()),
		UpdatedBy:     toNullString(username),
	}

	if req.School != "" {
		schoolUUID, err := uuid.Parse(req.School)
		if err != nil {
			return err
		}
		vehicle.SchoolUUID = &schoolUUID
	} else {
		vehicle.SchoolUUID = nil
	}

	isExistingVehicleNumber, err := service.vehicleRepository.CheckVehicleNumberExists(id, vehicle.VehicleNumber)
	if err != nil {
		return err
	}

	if isExistingVehicleNumber {
		return errors.New("Vehicle number already exists", 400)
	}

	err = service.vehicleRepository.UpdateVehicle(vehicle)
	if err != nil {
		return err
	}

	return nil
}

func (service *VehicleService) DeleteVehicle(id string, username string) error {
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	vehicle := entity.Vehicle{
		UUID:      parsedUUID,
		DeletedAt: toNullTime(time.Now()),
		DeletedBy: toNullString(username),
	}

	err = service.vehicleRepository.DeleteVehicle(vehicle)
	if err != nil {
		return err
	}

	return nil
}
