package dto

// type RoutePointRequestDTO struct {
// 	PointName string  `json:"point_name" validate:"required"`
// 	Order     int     `json:"point_order" validate:"required"`
// 	Latitude  float64 `json:"latitude" validate:"required"`
// 	Longitude float64 `json:"longitude" validate:"required"`
// }

type RouteRequestDTO struct {
	DriverUUID		 string					`json:"driver_uuid" validate:"required"`
	StudentUUID		 string					`json:"student_uuid" validate:"required"`
	SchoolUUID		 string					`json:"school_uuid" validate:"required"`
	RouteName        string                 `json:"route_name" validate:"required"`
	RouteDescription string                 `json:"route_description" validate:"required"`
}
