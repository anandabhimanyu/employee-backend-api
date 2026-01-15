package employee

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Handler struct {
	repo *Repository
}

func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// POST /employees
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(e)
}

// GET /employees/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	e, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(e)
}

// GET /employees
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	employees, err := h.repo.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employees)
}

// PUT /employees/{id}
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	var e Employee
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	e.ID = id

	if err := h.repo.Update(&e); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(e)
}

// DELETE /employees/{id}
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid employee id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
