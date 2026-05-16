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

type ProductHandler struct {
	repo *repositories.ProductRepository
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{repo: repositories.NewProductRepository()}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, 400, "payload tidak valid")
		return
	}
	p.ID = primitive.NewObjectID()
	if err := h.repo.Create(p); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 201, p)
}

func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAll()
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, products)
}

func (h *ProductHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	id, err := primitive.ObjectIDFromHex(mux.Vars(r)["id"])
	if err != nil {
		respondError(w, 400, "format ID tidak valid")
		return
	}
	p, err := h.repo.GetByID(id)
	if err != nil {
		respondError(w, 404, err.Error())
		return
	}
	respondJSON(w, 200, p)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var p models.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respondError(w, 400, "payload tidak valid")
		return
	}
	if err := h.repo.Update(id, p); err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.repo.Delete(mux.Vars(r)["id"])
	respondJSON(w, 200, map[string]string{"message": "product deleted"})
}

func (h *ProductHandler) Search(w http.ResponseWriter, r *http.Request) {
	var filter bson.M
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		respondError(w, 400, "filter tidak valid")
		return
	}
	products, err := h.repo.Search(filter)
	if err != nil {
		respondError(w, 500, err.Error())
		return
	}
	respondJSON(w, 200, products)
}
