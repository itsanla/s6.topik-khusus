package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"langstats-indonesia/data"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/api/languages", withCORS(handleAll))
	mux.HandleFunc("/api/languages/", withCORS(handleByName))
	mux.HandleFunc("/api/compare", withCORS(handleCompare))
	mux.HandleFunc("/api/stats", withCORS(handleStats))

	log.Printf("[LangStats] Server berjalan di port %s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func respond(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, map[string]string{
		"status":    "ok",
		"service":   "langstats-indonesia",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

func handleAll(w http.ResponseWriter, r *http.Request) {
	respond(w, http.StatusOK, map[string]any{
		"status": "ok",
		"total":  len(data.Languages),
		"data":   data.Languages,
	})
}

func handleByName(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/api/languages/")
	for _, lang := range data.Languages {
		if strings.EqualFold(lang.Name, name) {
			respond(w, http.StatusOK, lang)
			return
		}
	}
	respond(w, http.StatusNotFound, map[string]string{"error": "bahasa tidak ditemukan"})
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	type Row struct {
		Aspect string            `json:"aspect"`
		Values map[string]string `json:"values"`
	}
	rows := []Row{}
	aspects := []struct {
		label string
		fn    func(data.Language) string
	}{
		{"Paradigm", func(l data.Language) string { return l.Paradigm }},
		{"Performance", func(l data.Language) string { return l.Performance }},
		{"Concurrency", func(l data.Language) string { return l.Concurrency }},
		{"Compilation", func(l data.Language) string { return l.Compilation }},
		{"Error Model", func(l data.Language) string { return l.ErrorModel }},
	}
	for _, a := range aspects {
		row := Row{Aspect: a.label, Values: map[string]string{}}
		for _, lang := range data.Languages {
			row.Values[lang.Name] = a.fn(lang)
		}
		rows = append(rows, row)
	}
	respond(w, http.StatusOK, map[string]any{"comparison": rows})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	type Summary struct {
		Name       string `json:"name"`
		PopRank    int    `json:"popularity_rank"`
		LearnScore int    `json:"ease_of_learning"`
		SpeedScore int    `json:"speed"`
		EcoScore   int    `json:"ecosystem"`
	}
	summaries := make([]Summary, len(data.Languages))
	for i, l := range data.Languages {
		summaries[i] = Summary{l.Name, l.PopRank, l.LearnScore, l.SpeedScore, l.EcoScore}
	}
	respond(w, http.StatusOK, map[string]any{"stats": summaries})
}
