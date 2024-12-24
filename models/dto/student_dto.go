package dto

import ()

type StudentResponseDTO struct {
	UUID       string `json:"student_uuid"`
	FirstName  string `json:"student_first_name"`
	LastName   string `json:"student_last_name"`
	Gender     string `json:"student_gender" validate:"required,max=50"`
	Grade      string `json:"student_grade" validate:"required,max=50"`
	ParentUUID string `json:"parent_uuid,omitempty"`
	SchoolUUID string `json:"school_uuid"`
	SchoolName string `json:"school_name,omitempty"`
	StudentAddress   string `json:"student_address"`
	PickupPoint      string `json:"student_pickup_point"`
	ShuttleStatus string `json:"shuttle_status,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	CreatedBy  string `json:"created_by,omitempty"`
	UpdatedAt  string `json:"updated_at,omitempty"`
	UpdatedBy  string `json:"updated_by,omitempty"`
}


type StudentRequestDTO struct {
	StudentFirstName string `json:"student_first_name" validate:"required"`
	StudentLastName  string `json:"student_last_name" validate:"required"`
	StudentGender    Gender `json:"student_gender" validate:"required"`
	StudentGrade     string `json:"student_grade" validate:"required"`
	StudentAddress   string `json:"student_address" validate:"required"` // Menambahkan field student_address
	StudentPickupPoint map[string]float64 `json:"student_pickup_point" validate:"required"`
}

type StudentRequestByParentDTO struct {
	StudentFirstName string `json:"student_first_name" validate:"required"`
	StudentLastName  string `json:"student_last_name" validate:"required"`
	StudentGender    Gender `json:"student_gender" validate:"required"`
	StudentAddress   string `json:"student_address" validate:"required"` // Menambahkan field student_address
	StudentPickupPoint map[string]float64 `json:"student_pickup_point" validate:"required"`
}

type SchoolStudentParentRequestDTO struct {
	Student StudentRequestDTO `json:"student" validate:"required"`
	Parent  UserRequestsDTO   `json:"parent" validate:"required"`
}

type SchoolStudentParentResponseDTO struct {
	StudentUUID      string `json:"student_uuid"`
	ParentUUID       string `json:"parent_uuid,omitempty"`
	ParentName       string `json:"parent_name"`
	ParentPhone      string `json:"parent_phone"`
	StudentFirstName string `json:"student_first_name"`
	StudentLastName  string `json:"student_last_name"`
	StudentGender    Gender `json:"student_gender"`
	StudentGrade     string `json:"student_grade"`
	Address          string `json:"student_address"`
	PickupPoint      string `json:"student_pickup_point"` // Menambahkan field pickup_point
	ShuttleStatus      string `json:"shuttle_status"` // Menambahkan field pickup_point
	CreatedAt        string `json:"created_at,omitempty"`
	CreatedBy        string `json:"created_by,omitempty"`
	UpdatedAt        string `json:"updated_at,omitempty"`
	UpdatedBy        string `json:"updated_by,omitempty"`
}

