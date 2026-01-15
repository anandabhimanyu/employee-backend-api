package employee

import "time"

type Employee struct {
	ID        int       `json:"id"`
	FullName  string    `json:"full_name"`
	JobTitle  string    `json:"job_title"`
	Country   string    `json:"country"`
	Salary    float64   `json:"salary"`
	CreatedAt time.Time `json:"created_at"`
}
