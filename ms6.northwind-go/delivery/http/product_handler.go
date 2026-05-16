package http

import (
	"encoding/json"
	"net/http"
	"strings"

	"northwind-go/domain"
)

type ProductHandler struct {
	Usecase domain.ProductUsecase
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.TrimPrefix(r.URL.Path, "/products")
	path = strings.Trim(path, "/")

	switch {
	case r.URL.Path == "/products" && r.Method == http.MethodGet:
		h.fetchActive(w, r)
	case r.URL.Path == "/products" && r.Method == http.MethodPost:
		h.create(w, r)
	case strings.HasPrefix(r.URL.Path, "/products/") && r.Method == http.MethodGet:
		h.getOne(w, r, path)
	case strings.HasPrefix(r.URL.Path, "/products/price") && r.Method == http.MethodPut:
		h.updatePrice(w, r)
	case strings.HasPrefix(r.URL.Path, "/products/") && r.Method == http.MethodDelete:
		h.discontinue(w, r, path)
	default:
		respond(w, 405, map[string]string{"error": "method tidak diizinkan"})
	}
}

func (h *ProductHandler) fetchActive(w http.ResponseWriter, r *http.Request) {
	products, err := h.Usecase.GetActiveProducts(r.Context())
	if err != nil {
		respond(w, 500, map[string]string{"error": err.Error()})
		return
	}
	respond(w, 200, map[string]any{"total": len(products), "data": products})
}

func (h *ProductHandler) getOne(w http.ResponseWriter, r *http.Request, code string) {
	p, err := h.Usecase.GetProduct(r.Context(), code)
	if err != nil {
		respond(w, 404, map[string]string{"error": err.Error()})
		return
	}
	respond(w, 200, p)
}

func (h *ProductHandler) create(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		respond(w, 400, map[string]string{"error": "payload tidak valid"})
		return
	}
	if err := h.Usecase.CreateProduct(r.Context(), &p); err != nil {
		respond(w, 500, map[string]string{"error": err.Error()})
		return
	}
	respond(w, 201, p)
}

func (h *ProductHandler) updatePrice(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ProductCode string  `json:"product_code"`
		ListPrice   float64 `json:"list_price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		respond(w, 400, map[string]string{"error": "payload tidak valid"})
		return
	}
	if err := h.Usecase.UpdateProductPrice(r.Context(), payload.ProductCode, payload.ListPrice); err != nil {
		respond(w, 500, map[string]string{"error": err.Error()})
		return
	}
	respond(w, 200, map[string]string{"message": "harga produk berhasil diperbarui"})
}

func (h *ProductHandler) discontinue(w http.ResponseWriter, r *http.Request, code string) {
	if err := h.Usecase.DiscontinueProduct(r.Context(), code); err != nil {
		respond(w, 500, map[string]string{"error": err.Error()})
		return
	}
	respond(w, 200, map[string]string{"message": "produk berhasil di-discontinue"})
}

func respond(w http.ResponseWriter, code int, v any) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
