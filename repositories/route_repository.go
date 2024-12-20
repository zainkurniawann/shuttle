package repositories

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"errors"
// 	"time"

// 	"github.com/google/uuid"
// 	"shuttle/models/entity"
// )

// type RouteRepositoryInterface interface {
// 	AddRoute(route entity.Route) error
// 	GetRouteByID(routeID int64) (*entity.Route, error)
// }

// type RouteRepository struct {
// 	db *sql.DB
// }

// func NewRouteRepository(db *sql.DB) *RouteRepository {
// 	return &RouteRepository{db: db}
// }

// func (r *RouteRepository) AddRoute(route entity.Route) error {
// 	// Serialize RoutePoints ke JSON
// 	routePointsJSON, err := json.Marshal(entity.Route)
// 	if err != nil {
// 		return err
// 	}

// 	query := `
// 		INSERT INTO routes (
// 			route_uuid, school_uuid, route_name, route_description, route_status, 
// 			route_points, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
// 		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
// 	`

// 	_, err = r.db.Exec(query, route.RouteUUID, route.SchoolUUID, route.RouteName, route.RouteDescription,
// 		route.RouteStatus, routePointsJSON, time.Now(), route.CreatedBy, time.Now(), route.UpdatedBy, route.DeletedAt, route.DeletedBy)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

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
