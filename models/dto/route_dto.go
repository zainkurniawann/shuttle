package dto

import "database/sql"

///////////// ROUTE ASSIGNMENT /////////////
type RouteAssignmentDTO struct {
    StudentUUID      string `json:"student_uuid"`
    RouteName        string `json:"route_name"`
    RouteDescription string `json:"route_description"`
}

type RouteAssignmentResponseDTO struct {
    DriverUUID string             `json:"driver_uuid"`
    Students   []RouteAssignmentDTO  `json:"students"`
}

/////////// ROUTES //////////////////////
type RoutesResponseDTO struct {
	RouteNameUUID 			string       `json:"route_name_uuid,omitempty"`
	RouteName 				string       `json:"route_name,omitempty"`
	RouteDescription		string       `json:"route_description,omitempty"`
	CreatedAt        		string       `json:"created_at,omitempty"`
	CreatedBy        		string       `json:"created_by,omitempty"`
	UpdatedAt        		string       `json:"updated_at,omitempty"`
	UpdatedBy        		string       `json:"updated_by,omitempty"`
}

type RoutesRequestDTO struct {
	RouteName 				string       `json:"route_name" validate:"required"`
	RouteDescription		string       `json:"route_description" validate:"required"`
	CreatedAt        		string       `json:"created_at,omitempty"`
	CreatedBy        		string       `json:"created_by,omitempty"`
	UpdatedAt        		string       `json:"updated_at,omitempty"`
	UpdatedBy        		string       `json:"updated_by,omitempty"`
}

type RouteResponseByDriverDTO struct {
	RouteUUID          string         `json:"route_uuid,omitempty" db:"route_uuid"`
	StudentUUID        string         `json:"student_uuid,omitempty" db:"student_uuid"`
	DriverUUID         string         `json:"driver_uuid,omitempty" db:"driver_uuid"`
	SchoolUUID         string         `json:"school_uuid,omitempty" db:"school_uuid"`
	StudentFirstName   string         `json:"student_first_name,omitempty" db:"student_first_name"`
	StudentLastName    string         `json:"student_last_name,omitempty" db:"student_last_name"`
	StudentAddress     string         `json:"student_address,omitempty" db:"student_address"`
	StudentPickupPoint string         `json:"student_pickup_point,omitempty" db:"student_pickup_point"`
	ShuttleUUID        sql.NullString `db:"shuttle_uuid" json:"shuttle_uuid"`
	ShuttleStatus      sql.NullString `db:"shuttle_status" json:"shuttle_status"`
	SchoolName         string         `json:"school_name,omitempty" db:"school_name"`
	SchoolPoint        string         `json:"school_point,omitempty" db:"school_point"`
}