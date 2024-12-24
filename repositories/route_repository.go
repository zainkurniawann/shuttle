package repositories

import (
	// "database/sql"
	// "encoding/json"
	// "errors"
	// "time"
	"fmt"

	// "shuttle/models/dto"
	"shuttle/models/entity"

	"github.com/jmoiron/sqlx"
	// "github.com/google/uuid"
)

type RouteRepositoryInterface interface {
	AddRoute(route entity.Route) error
	UpdateRoute(routeUUID string, route entity.Route) error
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
