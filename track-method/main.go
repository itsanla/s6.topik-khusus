package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"track-method/config"
	"track-method/handler"
	"track-method/middleware"
	"track-method/repository"
	"track-method/usecase"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg := config.Load()

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Gagal terhubung ke Redis pada %s: %v", cfg.RedisAddr, err)
	}
	log.Printf("Redis terhubung pada %s", cfg.RedisAddr)

	trackRepo := repository.NewRedisTrackRepository(redisClient)
	trackUC := usecase.NewTrackUsecase(trackRepo)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer pingCancel()

		redisStatus := "ok"
		if err := redisClient.Ping(pingCtx).Err(); err != nil {
			redisStatus = fmt.Sprintf("error: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "track-method",
			"redis":   redisStatus,
		})
	})

	handler.NewTrackHandler(r, trackUC)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("Server berjalan di port %s (env: %s)", cfg.Port, cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server gagal berjalan: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Mematikan server dengan graceful shutdown...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server gagal shutdown: %v", err)
	}

	if err := redisClient.Close(); err != nil {
		log.Printf("Gagal menutup koneksi Redis: %v", err)
	}

	log.Println("Server berhasil dimatikan.")
}
