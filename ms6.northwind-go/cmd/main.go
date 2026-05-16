package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	deliveryHttp "northwind-go/delivery/http"
	"northwind-go/repository/mongodb"
	"northwind-go/usecase"

	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	port     := getEnv("PORT", "8081")
	dbName   := getEnv("MONGO_DB", "northwind")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mgo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("Gagal connect MongoDB: %v", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping MongoDB: %v", err)
	}
	log.Printf("[Northwind] Terhubung ke MongoDB: %s/%s", mongoURI, dbName)
	defer client.Disconnect(context.Background())

	db          := client.Database(dbName)
	productRepo := mongodb.NewMongoProductRepository(db)
	productUC   := usecase.NewProductUsecase(productRepo, 5*time.Second)

	handler := &deliveryHttp.ProductHandler{Usecase: productUC}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "ok",
			"service": "northwind-go",
		})
	})
	mux.Handle("/products", handler)
	mux.Handle("/products/", handler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	go func() {
		log.Printf("[Northwind] Server berjalan di port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("[Northwind] Graceful shutdown...")
	shutCtx, shutCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutCancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		fmt.Fprintf(os.Stderr, "Shutdown error: %v\n", err)
	}
}
