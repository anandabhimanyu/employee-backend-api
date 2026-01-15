package employee

import (
	"database/sql"
	"errors"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Create employee
func (r *Repository) Create(e *Employee) error {
	query := `
		INSERT INTO employees (full_name, job_title, country, salary)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRow(
		query,
		e.FullName,
		e.JobTitle,
		e.Country,
		e.Salary,
	).Scan(&e.ID, &e.CreatedAt)
}

// Get employee by ID
func (r *Repository) GetByID(id int) (*Employee, error) {
	query := `
		SELECT id, full_name, job_title, country, salary, created_at
		FROM employees
		WHERE id = $1
	`

	var e Employee
	err := r.db.QueryRow(query, id).Scan(
		&e.ID,
		&e.FullName,
		&e.JobTitle,
		&e.Country,
		&e.Salary,
		&e.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("employee not found")
	}

	return &e, err
}

// List all employees
func (r *Repository) List() ([]Employee, error) {
	query := `
		SELECT id, full_name, job_title, country, salary, created_at
		FROM employees
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var employees []Employee

	for rows.Next() {
		var e Employee
		if err := rows.Scan(
			&e.ID,
			&e.FullName,
			&e.JobTitle,
			&e.Country,
			&e.Salary,
			&e.CreatedAt,
		); err != nil {
			return nil, err
		}
		employees = append(employees, e)
	}

	return employees, nil
}

// Update employee
func (r *Repository) Update(e *Employee) error {
	query := `
		UPDATE employees
		SET full_name = $1,
		    job_title = $2,
		    country = $3,
		    salary = $4
		WHERE id = $5
	`

	res, err := r.db.Exec(
		query,
		e.FullName,
		e.JobTitle,
		e.Country,
		e.Salary,
		e.ID,
	)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("employee not found")
	}

	return nil
}

// Delete employee
func (r *Repository) Delete(id int) error {
	query := `DELETE FROM employees WHERE id = $1`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("employee not found")
	}

	return nil
}
