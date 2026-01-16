package employee

import (
	"database/sql"
	"errors"
	"strconv"
)

type Repository interface {
	Create(*Employee) error
	List(limit, offset int, country, sort, order string) ([]Employee, error)
	Count(country string) (int, error)
	GetByID(int) (*Employee, error)
	Update(*Employee) error
	Delete(int) error
}

type postgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(e *Employee) error {
	return r.db.QueryRow(`
		INSERT INTO employees (full_name, job_title, country, salary)
		VALUES ($1,$2,$3,$4)
		RETURNING id, created_at
	`,
		e.FullName, e.JobTitle, e.Country, e.Salary,
	).Scan(&e.ID, &e.CreatedAt)
}

func (r *postgresRepository) List(
	limit, offset int,
	country, sort, order string,
) ([]Employee, error) {

	// ✅ allow-list (VERY IMPORTANT)
	allowedSort := map[string]string{
		"id":         "id",
		"salary":     "salary",
		"created_at": "created_at",
	}

	if sort == "" {
		sort = "id"
	}
	if _, ok := allowedSort[sort]; !ok {
		sort = "id"
	}

	if order != "desc" {
		order = "asc"
	}

	query := `
		SELECT id, full_name, job_title, country, salary, created_at
		FROM employees
	`
	args := []interface{}{}
	idx := 1

	if country != "" {
		query += " WHERE country = $" + strconv.Itoa(idx)
		args = append(args, country)
		idx++
	}

	query += " ORDER BY " + allowedSort[sort] + " " + order

	if limit > 0 {
		query += " LIMIT $" + strconv.Itoa(idx)
		args = append(args, limit)
		idx++

		query += " OFFSET $" + strconv.Itoa(idx)
		args = append(args, offset)
	}

	rows, err := r.db.Query(query, args...)
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

func (r *postgresRepository) Count(country string) (int, error) {
	query := `SELECT COUNT(*) FROM employees`
	args := []interface{}{}

	if country != "" {
		query += " WHERE country = $1"
		args = append(args, country)
	}

	var total int
	err := r.db.QueryRow(query, args...).Scan(&total)
	return total, err
}

func (r *postgresRepository) GetByID(id int) (*Employee, error) {
	var e Employee
	err := r.db.QueryRow(`
		SELECT id, full_name, job_title, country, salary, created_at
		FROM employees WHERE id=$1
	`, id).Scan(
		&e.ID, &e.FullName, &e.JobTitle,
		&e.Country, &e.Salary, &e.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("employee not found")
	}
	return &e, err
}

func (r *postgresRepository) Update(e *Employee) error {
	res, err := r.db.Exec(`
		UPDATE employees
		SET full_name=$1, job_title=$2, country=$3, salary=$4
		WHERE id=$5
	`,
		e.FullName, e.JobTitle, e.Country, e.Salary, e.ID,
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

func (r *postgresRepository) Delete(id int) error {
	res, err := r.db.Exec(`DELETE FROM employees WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New("employee not found")
	}
	return nil
}
