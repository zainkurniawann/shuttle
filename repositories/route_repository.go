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

	// "shuttle/repositories"

	// "github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	// "github.com/google/uuid"
)

type RouteRepositoryInterface interface {
	FetchAllRoutes(schoolUUID string) ([]dto.RoutesResponseDTO, error)
	FetchSpecRoute(driverUUID, schoolUUID string) ([]entity.RouteAssignment, error)
	FetchAllRoutesByDriver(driverUUID string) ([]dto.RouteResponseByDriverDTO, error)
	FetchSpecRouteByDriver(driverUUID, studentUUID string) (*dto.RouteResponseByDriverDTO, error)
	GetRouteByStudentAndSchool(studentUUID, schoolUUID string) (*entity.RouteAssignment, error)
	ValidateDriverVehicle(driverUUID string) (bool, error)
	AddRoute(route entity.Routes) error
	GetSchoolUUIDByUserUUID(userUUID string, schoolUUID *string) error
	UpdateRoute(routeUUID string, route entity.RouteAssignment) error
	DeleteRoute(routeUUID string, username string) error
	// GetRouteByID(routeID int64) (*entity.RouteAssignment, error)
}

type routeRepository struct {
	DB *sqlx.DB
}

func NewRouteRepository(DB *sqlx.DB) *routeRepository {
	return &routeRepository{
		DB: DB,
	}
}

func (r *routeRepository) FetchAllRoutes(schoolUUID string) ([]dto.RoutesResponseDTO, error) {
	query := `
	SELECT 
		route_name_uuid, 
		route_name, 
		route_description, 
		created_at, 
		created_by, 
		updated_at, 
		updated_by
	FROM routes
	WHERE school_uuid = $1
	`

	rows, err := r.DB.Query(query, schoolUUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []dto.RoutesResponseDTO

	for rows.Next() {
		var route dto.RoutesResponseDTO
		var createdAt, updatedAt sql.NullTime
		var createdBy, updatedBy sql.NullString

		err := rows.Scan(
			&route.RouteNameUUID,
			&route.RouteName,
			&route.RouteDescription,
			&createdAt,
			&createdBy,
			&updatedAt,
			&updatedBy,
		)
		if err != nil {
			return nil, err
		}

		// Menangani nilai null
		if createdAt.Valid {
			route.CreatedAt = createdAt.Time.Format("2006-01-02 15:04:05")
		}
		if createdBy.Valid {
			route.CreatedBy = createdBy.String
		}
		if updatedAt.Valid {
			route.UpdatedAt = updatedAt.Time.Format("2006-01-02 15:04:05")
		}
		if updatedBy.Valid {
			route.UpdatedBy = updatedBy.String
		}

		routes = append(routes, route)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return routes, nil
}

func (repository *routeRepository) FetchSpecRoute(driverUUID, schoolUUID string) ([]entity.RouteAssignment, error) {
    log.Println("Starting FetchRoutesByDriverAndSchool repository")

    query := `
		SELECT 
			ra.driver_uuid,
			ra.student_uuid,
			r.route_name,
			r.route_description
		FROM route_assignment ra
		LEFT JOIN routes r 
			ON ra.route_name_uuid = r.route_name_uuid
		WHERE ra.driver_uuid = $1
		AND ra.school_uuid = $2
    `

    rows, err := repository.DB.Query(query, driverUUID, schoolUUID)
    if err != nil {
        log.Printf("Error querying routes: %v", err)
        return nil, fmt.Errorf("failed to fetch routes: %w", err)
    }
    defer rows.Close()

    var routes []entity.RouteAssignment
    for rows.Next() {
        var route entity.RouteAssignment
        if err := rows.Scan(&route.DriverUUID, &route.StudentUUID, &route.RouteName, &route.RouteDescription); err != nil {
            log.Printf("Error scanning route data: %v", err)
            return nil, fmt.Errorf("failed to scan route data: %w", err)
        }
        routes = append(routes, route)
    }

    return routes, nil
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
			st.shuttle_uuid,
			st.status AS shuttle_status,
			sc.school_name,
			sc.school_point
		FROM route_jawa r
		LEFT JOIN students s ON r.student_uuid = s.student_uuid
		LEFT JOIN schools sc ON r.school_uuid = sc.school_uuid
		LEFT JOIN shuttle st ON r.student_uuid = st.student_uuid AND DATE(st.created_at) = CURRENT_DATE
		WHERE r.driver_uuid = $1
		ORDER BY r.created_at ASC
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

func (r *routeRepository) GetRouteByStudentAndSchool(studentUUID, schoolUUID string) (*entity.RouteAssignment, error) {
	query := `
		SELECT route_id, route_uuid, driver_uuid, student_uuid, school_uuid, route_name, route_description, created_at, created_by
		FROM route_jawa
		WHERE student_uuid = $1 AND school_uuid = $2
		LIMIT 1
	`

	var route entity.RouteAssignment
	err := r.DB.QueryRow(query, studentUUID, schoolUUID).Scan(
		&route.RouteID, &route.RouteUUID, &route.DriverUUID, &route.StudentUUID, &route.SchoolUUID,
		&route.RouteName, &route.RouteDescription, &route.CreatedAt, &route.CreatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Tidak ditemukan, berarti route tidak ada
		}
		return nil, fmt.Errorf("failed to query database: %w", err)
	}

	return &route, nil
}

func (r *routeRepository) ValidateDriverVehicle(driverUUID string) (bool, error) {
	query := `
		SELECT 
			dd.vehicle_uuid, 
			v.driver_uuid
		FROM driver_details dd
		LEFT JOIN vehicles v ON dd.user_uuid = v.driver_uuid
		WHERE dd.user_uuid = $1
	`
	var vehicleUUID sql.NullString
	var driverUUIDFromVehicle sql.NullString

	err := r.DB.QueryRow(query, driverUUID).Scan(&vehicleUUID, &driverUUIDFromVehicle)
	if err != nil {
		return false, fmt.Errorf("failed to query driver details with vehicle join: %w", err)
	}

	// Pastikan vehicle_uuid dan driver_uuid tidak null
	if !vehicleUUID.Valid || !driverUUIDFromVehicle.Valid {
		return false, nil
	}

	// Jika kedua kolom valid (tidak kosong), kembalikan true
	return true, nil
}

func (r *routeRepository) AddRoute(route entity.Routes) error {
	query := `
		INSERT INTO routes (
			route_id, route_name_uuid, school_uuid, 
			route_name, route_description, created_at, created_by
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.DB.Exec(query,
		route.RouteID,
		route.RouteNameUUID.String(),
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

// Fungsi di userService untuk mengambil schoolUUID berdasarkan userUUID
func (r *routeRepository) GetSchoolUUIDByUserUUID(userUUID string, schoolUUID *string) error {
	query := `
		SELECT school_uuid
		FROM school_admin_details
		WHERE user_uuid = $1
	`
	err := r.DB.QueryRow(query, userUUID).Scan(schoolUUID)
	if err != nil {
		return fmt.Errorf("failed to get school UUID for user: %w", err)
	}
	return nil
}

func (r *routeRepository) UpdateRoute(routeUUID string, route entity.RouteAssignment) error {
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

// func (r *RouteRepository) GetRouteByID(routeID int64) (*entity.RouteAssignment, error) {
// 	query := `
// 		SELECT route_id, route_uuid, school_uuid, route_name, route_description,
// 		       route_status, route_points, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 		FROM routes WHERE route_id = ?
// 	`

// 	row := r.db.QueryRow(query, routeID)

// 	route := new(entity.RouteAssignment)
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
