package dto

// type RoutePointRequestDTO struct {
// 	PointName string  `json:"point_name" validate:"required"`
// 	Order     int     `json:"point_order" validate:"required"`
// 	Latitude  float64 `json:"latitude" validate:"required"`
// 	Longitude float64 `json:"longitude" validate:"required"`
// }

type RouteRequestDTO struct {
	DriverUUID		 string			`json:"driver_uuid" validate:"required"`
	StudentUUID		 string			`json:"student_uuid" validate:"required"`
	SchoolUUID		 string			`json:"school_uuid" validate:"required"`
	RouteName        string         `json:"route_name" validate:"required"`
	RouteDescription string         `json:"route_description" validate:"required"`
}

type RouteResponseDTO struct {
	RouteUUID		 string			`json:"route_uuid,omitempty"`
	DriverUUID		 string			`json:"driver_uuid,omitempty"`
	StudentUUID		 string			`json:"student_uuid,omitempty"`
	SchoolUUID		 string			`json:"school_uuid,omitempty"`
	RouteName        string         `json:"route_name,omitempty"`
	RouteDescription string         `json:"route_description,omitempty"`
	CreatedAt      	 string			`json:"created_at,omitempty"`
	CreatedBy      	 string 		`json:"created_by,omitempty"`
	UpdatedAt      	 string 		`json:"updated_at,omitempty"`
	UpdatedBy      	 string 		`json:"updated_by,omitempty"`
}

type RouteResponseByDriverDTO struct {
    RouteUUID          string `db:"route_uuid"`
    StudentUUID        string `db:"student_uuid"`
    DriverUUID         string `db:"driver_uuid"`
    SchoolUUID         string `db:"school_uuid"`
    StudentFirstName   string `db:"student_first_name"`
    StudentLastName    string `db:"student_last_name"`
    StudentAddress     string `db:"student_address"`
    StudentPickupPoint string `db:"student_pickup_point"`
    SchoolName         string `db:"school_name"`
    SchoolPoint        string `db:"school_point"`
}













// -- Table: public.route_jawa

// -- DROP TABLE IF EXISTS public.route_jawa;

// CREATE TABLE IF NOT EXISTS public.route_jawa
// (
//     route_id bigint NOT NULL,
//     route_uuid uuid NOT NULL,
//     driver_uuid uuid NOT NULL,
//     student_uuid uuid NOT NULL,
//     route_name character varying(100) COLLATE pg_catalog."default" NOT NULL,
//     route_description text COLLATE pg_catalog."default",
//     created_at timestamp with time zone NOT NULL,
//     created_by character varying(255) COLLATE pg_catalog."default" NOT NULL,
//     updated_at timestamp with time zone,
//     updated_by character varying(255) COLLATE pg_catalog."default",
//     deleted_at timestamp with time zone,
//     deleted_by character varying(255) COLLATE pg_catalog."default",
//     school_uuid uuid NOT NULL,
//     CONSTRAINT route_jawa_pkey PRIMARY KEY (route_id),
//     CONSTRAINT fk_driver_uuid FOREIGN KEY (driver_uuid)
//         REFERENCES public.users (user_uuid) MATCH SIMPLE
//         ON UPDATE NO ACTION
//         ON DELETE NO ACTION,
//     CONSTRAINT fk_school_uuid FOREIGN KEY (school_uuid)
//         REFERENCES public.schools (school_uuid) MATCH SIMPLE
//         ON UPDATE NO ACTION
//         ON DELETE NO ACTION,
//     CONSTRAINT fk_student_uuid FOREIGN KEY (student_uuid)
//         REFERENCES public.students (student_uuid) MATCH SIMPLE
//         ON UPDATE NO ACTION
//         ON DELETE NO ACTION
// )

// TABLESPACE pg_default;

// ALTER TABLE IF EXISTS public.route_jawa
//     OWNER to postgres;
