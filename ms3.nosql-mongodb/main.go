package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nosql-mongodb/internal/config"
	"nosql-mongodb/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {
	config.ConnectMySQL()
	config.ConnectMongoDB()
	defer config.CloseConnections()

	r := mux.NewRouter()

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": "nosql-mongodb"})
	}).Methods(http.MethodGet)

	ph := handlers.NewProductHandler()
	r.HandleFunc("/product", ph.Create).Methods(http.MethodPost)
	r.HandleFunc("/products", ph.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/product/{id}", ph.GetOne).Methods(http.MethodGet)
	r.HandleFunc("/product/{id}", ph.Update).Methods(http.MethodPut)
	r.HandleFunc("/product/{id}", ph.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/search", ph.Search).Methods(http.MethodPost)

	ch := handlers.NewCustomerHandler()
	r.HandleFunc("/customer", ch.Create).Methods(http.MethodPost)
	r.HandleFunc("/customers", ch.GetAll).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", ch.GetOne).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", ch.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/customer/search", ch.Search).Methods(http.MethodPost)

	port := os.Getenv("PORT")
	if port == "" {
		port = "9090"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("[NoSQL-MongoDB] Server berjalan di port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("[NoSQL-MongoDB] Server berhenti.")
}
