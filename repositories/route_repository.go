package repositories

import (
	// "database/sql"
	// "encoding/json"
	// "errors"
	// "time"
	"database/sql"
	"fmt"
	"log"

	// "shuttle/models/dto"
	"shuttle/models/dto"
	"shuttle/models/entity"

	"github.com/jmoiron/sqlx"
	// "github.com/google/uuid"
)

type RouteRepositoryInterface interface {
	FetchAllRoutes() ([]entity.Route, error)
	FetchSpecRoute(routeUUID string) (entity.Route, error)
	FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)
	FetchSpecRouteByDriver(driverUUID, studentUUID string) (*dto.RouteResponseByDriverDTO, error)
	AddRoute(route entity.Route) error
	UpdateRoute(routeUUID string, route entity.Route) error
	DeleteRoute(routeUUID string, username string) error
	// GetRouteByID(routeID int64) (*entity.Route, error)
}

type routeRepository struct {
	DB *sqlx.DB
}

func NewRouteRepository(DB *sqlx.DB) *routeRepository {
	return &routeRepository{
		DB: DB,
	}
}

func (r *routeRepository) FetchAllRoutes() ([]entity.Route, error) {
	// Query untuk mengambil semua data dari tabel routes
	query := `
		SELECT route_uuid, driver_uuid, student_uuid, school_uuid, 
		       route_name, route_description, created_at, created_by, updated_at, updated_by
		FROM route_jawa
	`
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Slice untuk menyimpan hasil query
	var routes []entity.Route

	// Iterasi melalui hasil query
	for rows.Next() {
		var route entity.Route
		if err := rows.Scan(
			&route.RouteUUID,
			&route.DriverUUID,
			&route.StudentUUID,
			&route.SchoolUUID,
			&route.RouteName,
			&route.RouteDescription,
			&route.CreatedAt,
			&route.CreatedBy,
			&route.UpdatedAt,
			&route.UpdatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		routes = append(routes, route)
	}

	// Cek jika ada error setelah iterasi
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over rows: %w", err)
	}

	return routes, nil
}

func (repository *routeRepository) FetchSpecRoute(routeUUID string) (entity.Route, error) {
	log.Println("Starting GetRouteByUUID repository")

	// Membuat query untuk mengambil route berdasarkan UUID
	var route entity.Route
	query := `SELECT route_uuid, driver_uuid, student_uuid, school_uuid, route_name, 
		route_description, created_at, created_by, updated_at, updated_by
		FROM route_jawa WHERE route_uuid = $1`

	// Menjalankan query
	err := repository.DB.QueryRow(query, routeUUID).Scan(
		&route.RouteUUID,
		&route.DriverUUID,
		&route.StudentUUID,
		&route.SchoolUUID,
		&route.RouteName,
		&route.RouteDescription,
		&route.CreatedAt,
		&route.CreatedBy,
		&route.UpdatedAt,
		&route.UpdatedBy,
	)
	if err != nil {
		log.Printf("Error querying route by UUID: %v", err)
		if err == sql.ErrNoRows {
			return route, fmt.Errorf("route not found")
		}
		return route, fmt.Errorf("failed to get route by UUID: %w", err)
	}

	log.Printf("Route found: %+v", route)
	return route, nil
}

func (repo *routeRepository) FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error) {
	query := `
		SELECT
			r.route_uuid,
			r.student_uuid,
			r.driver_uuid,
			r.school_uuid,
			s.student_first_name,
			s.student_last_name,
			s.student_address,
			s.student_pickup_point,
			sc.school_name,
			sc.school_point
		FROM route_jawa r
		LEFT JOIN students s ON r.student_uuid = s.student_uuid
		LEFT JOIN schools sc ON r.school_uuid = sc.school_uuid
		WHERE r.driver_uuid = $1
	`
	var routes []dto.RouteResponseByDriverDTO
	err := repo.DB.Select(&routes, query, driverUUID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	return routes, nil
}

func (repo *routeRepository) FetchSpecRouteByDriver(driverUUID, studentUUID string) (*dto.RouteResponseByDriverDTO, error) {
	query := `
		SELECT
			r.route_uuid,
			r.student_uuid,
			r.driver_uuid,
			r.school_uuid,
			s.student_first_name,
			s.student_last_name,
			s.student_address,
			s.student_pickup_point,
			sc.school_name,
			sc.school_point
		FROM route_jawa r
		LEFT JOIN students s ON r.student_uuid = s.student_uuid
		LEFT JOIN schools sc ON r.school_uuid = sc.school_uuid
		WHERE r.driver_uuid = $1 AND r.student_uuid = $2
	`
	var route dto.RouteResponseByDriverDTO
	err := repo.DB.Get(&route, query, driverUUID, studentUUID)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	return &route, nil
}

func (r *routeRepository) AddRoute(route entity.Route) error {
	query := `
		INSERT INTO route_jawa (
			route_id, route_uuid, driver_uuid, student_uuid, school_uuid, 
			route_name, route_description, created_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.DB.Exec(query,
		route.RouteID,
		route.RouteUUID.String(),
		route.DriverUUID.String(),
		route.StudentUUID.String(),
		route.SchoolUUID.String(),
		route.RouteName,
		route.RouteDescription,
		route.CreatedAt.Time,
		route.CreatedBy.String,
	)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	return nil
}

func (r *routeRepository) UpdateRoute(routeUUID string, route entity.Route) error {
	query := `
		UPDATE route_jawa
		SET 
			driver_uuid = $1,
			student_uuid = $2,
			school_uuid = $3,
			route_name = $4,
			route_description = $5,
			updated_at = $6,
			updated_by = $7
		WHERE 
			route_uuid = $8
	`

	_, err := r.DB.Exec(query,
		route.DriverUUID.String(),
		route.StudentUUID.String(),
		route.SchoolUUID.String(),
		route.RouteName,
		route.RouteDescription,
		route.UpdatedAt.Time,
		route.UpdatedBy.String,
		routeUUID,
	)
	if err != nil {
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

func (repo *routeRepository) DeleteRoute(routeUUID string, username string) error {
	// Query untuk menghapus route berdasarkan UUID
	query := `DELETE FROM route_jawa WHERE route_uuid = $1`
	_, err := repo.DB.Exec(query, routeUUID)
	if err != nil {
		return fmt.Errorf("failed to execute delete query: %w", err)
	}

	// Log penghapusan oleh username
	log.Printf("Route deleted by: %s", username)
	return nil
}









































// func (r *RouteRepository) GetRouteByID(routeID int64) (*entity.Route, error) {
// 	query := `
// 		SELECT route_id, route_uuid, school_uuid, route_name, route_description,
// 		       route_status, route_points, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 		FROM routes WHERE route_id = ?
// 	`

// 	row := r.db.QueryRow(query, routeID)

// 	route := new(entity.Route)
// 	var routePointsJSON []byte
// 	if err := row.Scan(
// 		&route.RouteID, &route.RouteUUID, &route.SchoolUUID, &route.RouteName, &route.RouteDescription,
// 		&route.RouteStatus, &routePointsJSON, &route.CreatedAt, &route.CreatedBy, &route.UpdatedAt, &route.UpdatedBy,
// 		&route.DeletedAt, &route.DeletedBy,
// 	); err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, errors.New("route not found")
// 		}
// 		return nil, err
// 	}

// 	// Deserialize RoutePoints dari JSON
// 	if err := json.Unmarshal(routePointsJSON, &route.RoutePoints); err != nil {
// 		return nil, err
// 	}

// 	return route, nil
// }
