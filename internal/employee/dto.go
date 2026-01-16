package employee

type CreateEmployeeRequest struct {
	FullName string  `json:"full_name" binding:"required,min=3"`
	JobTitle string  `json:"job_title" binding:"required"`
	Country  string  `json:"country" binding:"required"`
	Salary   float64 `json:"salary" binding:"required,gt=0"`
}

type UpdateEmployeeRequest struct {
	FullName string  `json:"full_name" binding:"required,min=3"`
	JobTitle string  `json:"job_title" binding:"required"`
	Country  string  `json:"country" binding:"required"`
	Salary   float64 `json:"salary" binding:"required,gt=0"`
}
