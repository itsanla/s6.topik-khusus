package handlers

import (
	"encoding/json"
	"net/http"

	"nosql-mongodb/internal/models"
	"nosql-mongodb/internal/repositories"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerHandler struct {
	repo *repositories.CustomerRepository
}

func NewCustomerHandler() *CustomerHandler {
	return &CustomerHandler{repo: repositories.NewCustomerRepository()}
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var c models.Customer
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		respondError(w, 400, "payload tidak valid")
		return
	}
	c.ID = primitive.NewObjectID()
	if err := h.repo.Create(c); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, c)
}

func (h *CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	customers, err := h.repo.GetAll()
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, customers)
}

func (h *CustomerHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, 400, "format ID tidak valid")
		return
	}
	c, err := h.repo.GetByID(id)
	if err != nil {
		respondError(w, 404, err.Error())
		return
	}
	respondJSON(w, 200, c)
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.repo.Delete(mux.Vars(r)["id"])
	respondJSON(w, 200, map[string]string{"message": "customer deleted"})
}

func (h *CustomerHandler) Search(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		respondError(w, 400, "filter tidak valid")
		return
	}
	customers, err := h.repo.Search(filter)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, customers)
}
