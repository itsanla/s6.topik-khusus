package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"elasticsearch-guide/store"
)

var db = store.New()

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/health", handleHealth)

	// Index operations: PUT /:index  DELETE /:index  HEAD /:index
	// Document operations: PUT /:index/_doc/:id  GET /:index/_doc/:id  DELETE /:index/_doc/:id
	// Search: POST /:index/_search
	// Update: POST /:index/_update/:id
	// Cluster info: GET /
	// List indices: GET /_cat/indices
	mux.HandleFunc("/", router)

	log.Printf("[Elasticsearch Guide] Simulator berjalan di port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func router(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")

	switch {
	case path == "":
		handleClusterInfo(w, r)
	case path == "_cat/indices":
		handleCatIndices(w, r)
	case len(parts) == 1:
		handleIndex(w, r, parts[0])
	case len(parts) == 3 && parts[1] == "_doc":
		handleDoc(w, r, parts[0], parts[2])
	case len(parts) == 2 && parts[1] == "_search":
		handleSearch(w, r, parts[0])
	case len(parts) == 3 && parts[1] == "_update":
		handleUpdate(w, r, parts[0], parts[2])
	case len(parts) == 2 && parts[1] == "_count":
		handleCount(w, r, parts[0])
	default:
		respond(w, 404, map[string]any{"error": "endpoint tidak dikenal"})
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respond(w, 200, map[string]any{
		"status":    "ok",
		"service":   "elasticsearch-guide",
		"simulator": true,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func handleClusterInfo(w http.ResponseWriter, r *http.Request) {
	respond(w, 200, map[string]any{
		"name":         "es-simulator",
		"cluster_name": "s6-topik-khusus",
		"version":      map[string]any{"number": "8.x-simulated"},
		"tagline":      "You Know, for Search (Simulator)",
	})
}

func handleCatIndices(w http.ResponseWriter, r *http.Request) {
	indices := db.ListIndices()
	result := make([]map[string]any, 0, len(indices))
	for _, idx := range indices {
		result = append(result, map[string]any{
			"index":      idx,
			"docs.count": db.Count(idx),
			"status":     "open",
		})
	}
	respond(w, 200, result)
}

func handleIndex(w http.ResponseWriter, r *http.Request, index string) {
	switch r.Method {
	case http.MethodPut:
		if db.CreateIndex(index) {
			respond(w, 200, map[string]any{"acknowledged": true, "index": index})
		} else {
			respond(w, 400, map[string]any{"error": fmt.Sprintf("index '%s' sudah ada", index)})
		}
	case http.MethodDelete:
		if db.DeleteIndex(index) {
			respond(w, 200, map[string]any{"acknowledged": true})
		} else {
			respond(w, 404, map[string]any{"error": fmt.Sprintf("index '%s' tidak ditemukan", index)})
		}
	case http.MethodHead:
		if db.IndexExists(index) {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(404)
		}
	default:
		respond(w, 405, map[string]any{"error": "method tidak diizinkan"})
	}
}

func handleDoc(w http.ResponseWriter, r *http.Request, index, id string) {
	switch r.Method {
	case http.MethodPut:
		var doc map[string]any
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			respond(w, 400, map[string]any{"error": "payload tidak valid"})
			return
		}
		created := db.PutDoc(index, id, doc)
		result := "updated"
		if created {
			result = "created"
		}
		respond(w, 200, map[string]any{"_index": index, "_id": id, "result": result})
	case http.MethodGet:
		doc, ok := db.GetDoc(index, id)
		if !ok {
			respond(w, 404, map[string]any{"found": false, "_index": index, "_id": id})
			return
		}
		respond(w, 200, map[string]any{"_index": index, "_id": id, "found": true, "_source": doc})
	case http.MethodDelete:
		if db.DeleteDoc(index, id) {
			respond(w, 200, map[string]any{"_index": index, "_id": id, "result": "deleted"})
		} else {
			respond(w, 404, map[string]any{"result": "not_found"})
		}
	default:
		respond(w, 405, map[string]any{"error": "method tidak diizinkan"})
	}
}

func handleSearch(w http.ResponseWriter, r *http.Request, index string) {
	var body map[string]any
	json.NewDecoder(r.Body).Decode(&body)

	query := ""
	if q, ok := body["query"].(map[string]any); ok {
		if ms, ok := q["match"].(map[string]any); ok {
			for _, v := range ms {
				query = fmt.Sprintf("%v", v)
				break
			}
		}
		if qs, ok := q["query_string"].(map[string]any); ok {
			query = fmt.Sprintf("%v", qs["query"])
		}
	}

	hits := db.Search(index, query)
	respond(w, 200, map[string]any{
		"hits": map[string]any{
			"total": map[string]any{"value": len(hits)},
			"hits":  hits,
		},
	})
}

func handleUpdate(w http.ResponseWriter, r *http.Request, index, id string) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		respond(w, 400, map[string]any{"error": "payload tidak valid"})
		return
	}
	doc, _ := body["doc"].(map[string]any)
	if err := db.UpdateDoc(index, id, doc); err != nil {
		respond(w, 404, map[string]any{"error": err.Error()})
		return
	}
	respond(w, 200, map[string]any{"_index": index, "_id": id, "result": "updated"})
}

func handleCount(w http.ResponseWriter, r *http.Request, index string) {
	respond(w, 200, map[string]any{"count": db.Count(index), "_index": index})
}

func respond(w http.ResponseWriter, code int, v any) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
