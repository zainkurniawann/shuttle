package repositories

import (
	"log"
	"shuttle/models/entity"

	"github.com/jmoiron/sqlx"
)

type ChildernRepositoryInterface interface {
	FetchAllChilderns(id string) ([]entity.Student, error)
	FetchSpecChildern(id string) (entity.Student, error)
	UpdateChildern(tx *sqlx.Tx, student entity.Student, studentUUID string) error
}

type childernRepository struct {
	DB *sqlx.DB
}

func NewChildernRepository(DB *sqlx.DB) ChildernRepositoryInterface {
	return &childernRepository{
		DB: DB,
	}
}

func (repositories *childernRepository) FetchAllChilderns(id string) ([]entity.Student, error) {
	var childerns []entity.Student

	query := `
        SELECT 
        s.student_uuid,
        s.student_first_name,
        s.student_last_name,
        s.student_gender,
        s.student_grade,
        s.parent_uuid,
        s.school_uuid,
        st.status AS shuttle_status
    FROM students s
    LEFT JOIN shuttle st 
        ON s.student_uuid = st.student_uuid
    WHERE s.parent_uuid = $1
    `

	rows, err := repositories.DB.Queryx(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var childern entity.Student

		if err := rows.Scan(
			&childern.UUID,
			&childern.FirstName,
			&childern.LastName,
			&childern.Gender,
			&childern.Grade,
			&childern.ParentUUID,
			&childern.SchoolUUID,
			&childern.ShuttleStatus, // Tambahkan kolom ini
		); err != nil {
			return nil, err
		}

		childerns = append(childerns, childern)
	}

	return childerns, nil
}

func (repositories *childernRepository) FetchSpecChildern(id string) (entity.Student, error) {
	var childern entity.Student

	query := `
		SELECT 
            s.student_uuid,
            s.student_first_name,
            s.student_last_name,
            s.student_gender,
            s.student_grade,
            s.parent_uuid,
            s.school_uuid,
            sc.school_name
        FROM students s
        JOIN schools sc ON s.school_uuid = sc.school_uuid
        WHERE s.student_uuid = $1
	`

	err := repositories.DB.QueryRowx(query, id).Scan(
		&childern.UUID,
		&childern.FirstName,
		&childern.LastName,
		&childern.Gender,
		&childern.Grade,
		&childern.ParentUUID,
		&childern.SchoolUUID,
	)
	log.Println("SchoolUUID:", childern.SchoolUUID)

	if err != nil {
		return entity.Student{}, err
	}

	return childern, nil
}

func (r *childernRepository) UpdateChildern(tx *sqlx.Tx, student entity.Student, studentUUID string) error {
	query := `
        UPDATE students
        SET student_first_name = $1, student_last_name = $2, student_gender = $3, updated_at = NOW(), updated_by = $4
        WHERE student_uuid = $5`
	_, err := tx.Exec(query, student.FirstName, student.LastName, student.Gender, student.UpdatedBy, studentUUID)
	return err
}
