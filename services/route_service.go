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
	AddRoute(route dto.RouteRequestDTO, schoolUUID, username string) error
	UpdateRoute(routeUUID string, route dto.RouteRequestDTO, username string) error
}

type routeService struct {
	routeRepository repositories.RouteRepositoryInterface
}

func NewRouteService(routeRepository repositories.RouteRepositoryInterface) RouteServiceInterface {
	return &routeService{
		routeRepository: routeRepository,
	}
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
