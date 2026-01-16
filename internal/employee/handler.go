package employee

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// ================= CREATE =================
func (h *Handler) Create(c *gin.Context) {
	var e Employee

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid JSON payload",
		})
		return
	}

	// ✅ VALIDATION
	if e.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "full_name is required"})
		return
	}
	if e.JobTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "job_title is required"})
		return
	}
	if e.Country == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "country is required"})
		return
	}
	if e.Salary <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "salary must be greater than 0"})
		return
	}

	if err := h.repo.Create(&e); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, e)
}

// ================= LIST =================
func (h *Handler) List(c *gin.Context) {
	// 1️⃣ Read query params
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	country := c.Query("country")

	// 2️⃣ Fetch data
	list, err := h.repo.List(limit, offset, country)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 3️⃣ Return structured response
	c.JSON(http.StatusOK, gin.H{
		"data": list,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(list),
		},
	})
}

// ================= GET BY ID =================
func (h *Handler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid employee id",
		})
		return
	}

	e, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, e)
}

// ================= UPDATE =================
func (h *Handler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	var e Employee
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON payload"})
		return
	}

	// ✅ VALIDATION
	if e.FullName == "" || e.JobTitle == "" || e.Country == "" || e.Salary <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	e.ID = id

	if err := h.repo.Update(&e); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee updated"})
}

// ================= DELETE =================
func (h *Handler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "employee deleted"})
}
