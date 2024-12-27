package services

import (
	"database/sql"
	"fmt"
	"log"
	"shuttle/errors"
	"shuttle/models/dto"
	"shuttle/models/entity"
	"shuttle/repositories"
	"time"

	"github.com/google/uuid"
)

type RouteServiceInterface interface {
	GetAllRoutes() ([]dto.RouteResponseDTO, error)
	GetSpecRoute(routeUUID string) (dto.RouteResponseDTO, error)
	GetAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)
	GetSpecRouteByDriver(driverUUID, studentUUID string) (*dto.RouteResponseByDriverDTO, error)
	AddRoute(route dto.RouteRequestDTO, schoolUUID, username string) error
	UpdateRoute(routeUUID string, route dto.RouteRequestDTO, username string) error
	DeleteRoute(routeUUID string, username string) error
}

type routeService struct {
	routeRepository repositories.RouteRepositoryInterface
}

func NewRouteService(routeRepository repositories.RouteRepositoryInterface) RouteServiceInterface {
	return &routeService{
		routeRepository: routeRepository,
	}
}

func (service *routeService) GetAllRoutes() ([]dto.RouteResponseDTO, error) {
	log.Println("Starting GetAllRoutes service")

	// Memanggil repository untuk mengambil semua routes
	routes, err := service.routeRepository.FetchAllRoutes()
	if err != nil {
		log.Printf("Failed to get routes: %v", err)
		return nil, fmt.Errorf("failed to get routes: %w", err)
	}

	// Mapping hasil ke DTO
	var routeDTOs []dto.RouteResponseDTO
	for _, route := range routes {
		// Mengonversi time.Time ke string
		createdAt := route.CreatedAt.Time.Format("2006-01-02 15:04:05")
		
		// Cek apakah UpdatedAt valid
		var updatedAt string
		if route.UpdatedAt.Valid {
			updatedAt = route.UpdatedAt.Time.Format("2006-01-02 15:04:05")
		} else {
			updatedAt = "N/A" // Atau bisa juga ""
		}

		// Cek apakah UpdatedBy valid
		var updatedBy string
		if route.UpdatedBy.Valid {
			updatedBy = route.UpdatedBy.String
		} else {
			updatedBy = "N/A" // Atau bisa juga ""
		}

		// Menambahkan UserUsername dan StudentName
		routeDTOs = append(routeDTOs, dto.RouteResponseDTO{
			RouteUUID:        route.RouteUUID.String(),
			DriverUUID:       route.DriverUUID.String(),
			UserUsername:     route.UserUsername, // Menambahkan UserUsername
			StudentUUID:      route.StudentUUID.String(),
			StudentName:      route.StudentName,   // Menambahkan StudentName
			SchoolUUID:       route.SchoolUUID.String(),
			RouteName:        route.RouteName,
			RouteDescription: route.RouteDescription,
			CreatedAt:        createdAt, // Menggunakan string
			CreatedBy:        route.CreatedBy.String,
			UpdatedAt:        updatedAt, // Menggunakan string atau "N/A"
			UpdatedBy:        updatedBy, // Menggunakan string atau "N/A"
		})
	}

	log.Println("Routes fetched successfully")
	return routeDTOs, nil
}

func (service *routeService) GetSpecRoute(routeUUID string) (dto.RouteResponseDTO, error) {
	log.Println("Starting GetRouteByUUID service")

	// Memanggil repository untuk mendapatkan route berdasarkan UUID
	route, err := service.routeRepository.FetchSpecRoute(routeUUID)
	if err != nil {
		log.Printf("Failed to get route from repository: %v", err)
		return dto.RouteResponseDTO{}, fmt.Errorf("failed to get route: %w", err)
	}

	// Mengonversi waktu menjadi string dan mengatur nilai untuk UpdatedAt dan UpdatedBy
	createdAt := route.CreatedAt.Time.Format("2006-01-02 15:04:05")

	var updatedAt string
	if route.UpdatedAt.Valid {
		updatedAt = route.UpdatedAt.Time.Format("2006-01-02 15:04:05")
	} else {
		updatedAt = "N/A"
	}

	var updatedBy string
	if route.UpdatedBy.Valid {
		updatedBy = route.UpdatedBy.String
	} else {
		updatedBy = "N/A"
	}

	// Return RouteResponseDTO yang sudah di-mapping
	return dto.RouteResponseDTO{
		RouteUUID:        route.RouteUUID.String(),
		DriverUUID:       route.DriverUUID.String(),
		StudentUUID:      route.StudentUUID.String(),
		SchoolUUID:       route.SchoolUUID.String(),
		RouteName:        route.RouteName,
		RouteDescription: route.RouteDescription,
		CreatedAt:        createdAt,
		CreatedBy:        route.CreatedBy.String,
		UpdatedAt:        updatedAt,
		UpdatedBy:        updatedBy,
	}, nil
}

func (service *routeService) GetAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error) {
	routes, err := service.routeRepository.FetchAllRoutesByDriver(driverUUID)
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (service *routeService) GetSpecRouteByDriver(driverUUID, studentUUID string) (*dto.RouteResponseByDriverDTO, error) {
	route, err := service.routeRepository.FetchSpecRouteByDriver(driverUUID, studentUUID)
	if err != nil {
		return nil, err
	}
	return route, nil
}


func (service *routeService) AddRoute(route dto.RouteRequestDTO, schoolUUID, username string) error {
	log.Println("Starting AddRoute service")

	// Mapping DTO ke Entity
	routeEntity := entity.Route{
		RouteID:          time.Now().UnixMilli()*1e6 + int64(uuid.New().ID()%1e6),
		RouteUUID:        uuid.New(),
		DriverUUID:       uuid.MustParse(route.DriverUUID),
		StudentUUID:      uuid.MustParse(route.StudentUUID),
		SchoolUUID:       uuid.MustParse(schoolUUID),
		RouteName:        route.RouteName,
		RouteDescription: route.RouteDescription,
		CreatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
		CreatedBy:        sql.NullString{String: username, Valid: true},
	}

	// Log entitas route yang akan disimpan
	log.Printf("Route entity to be saved: %+v", routeEntity)

	// Simpan ke repository
	log.Println("Calling repository to add route")
	if err := service.routeRepository.AddRoute(routeEntity); err != nil {
		log.Printf("Failed to add route: %v", err)
		return errors.New("Failed to add route", 500)
	}

	log.Println("Route added successfully")
	return nil
}

func (service *routeService) UpdateRoute(routeUUID string, route dto.RouteRequestDTO, username string) error {
	log.Println("Starting UpdateRoute service")

	// Validasi dan mapping DTO ke Entity
	log.Println("Mapping DTO to Entity")
	routeEntity := entity.Route{
		DriverUUID:       uuid.MustParse(route.DriverUUID),
		StudentUUID:      uuid.MustParse(route.StudentUUID),
		SchoolUUID:       uuid.MustParse(route.SchoolUUID),
		RouteName:        route.RouteName,
		RouteDescription: route.RouteDescription,
		UpdatedAt:        sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedBy:        sql.NullString{String: username, Valid: true}, // username dari context
	}

	// Log detail entity yang akan di-update
	log.Printf("Route UUID to be updated: %s", routeUUID)
	log.Printf("Route entity data to be updated: %+v", routeEntity)

	// Panggil repository untuk melakukan update
	log.Println("Calling repository to update route")
	if err := service.routeRepository.UpdateRoute(routeUUID, routeEntity); err != nil {
		log.Printf("Failed to update route: %v", err)
		return fmt.Errorf("failed to update route: %w", err)
	}

	log.Println("Route updated successfully")
	return nil
}

func (service *routeService) DeleteRoute(routeUUID string, username string) error {
	log.Println("Starting DeleteRoute service")

	// Log informasi yang akan digunakan untuk delete
	log.Printf("Route UUID to be deleted: %s", routeUUID)

	// Panggil repository untuk delete route
	log.Println("Calling repository to delete route")
	if err := service.routeRepository.DeleteRoute(routeUUID, username); err != nil {
		log.Printf("Failed to delete route: %v", err)
		return fmt.Errorf("failed to delete route: %w", err)
	}

	log.Println("Route deleted successfully")
	return nil
}
