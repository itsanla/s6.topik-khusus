package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *sql.DB
var MongoClient *mongo.Client

const (
	MongoDatabaseName       = "store"
	MongoProductCollection  = "products"
	MongoCustomerCollection = "customers"
)

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func ConnectMySQL() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("MYSQL_USER", "root"),
		getEnv("MYSQL_PASS", ""),
		getEnv("MYSQL_HOST", "127.0.0.1"),
		getEnv("MYSQL_PORT", "3306"),
		getEnv("MYSQL_DB", "db_store"),
	)
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Gagal buka koneksi MySQL: %v", err)
	}
	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(3 * time.Minute)
	if err = DB.Ping(); err != nil {
		log.Fatalf("Gagal ping MySQL: %v", err)
	}
	log.Printf("[Config] MySQL terhubung: %s/%s", getEnv("MYSQL_HOST", "127.0.0.1"), getEnv("MYSQL_DB", "db_store"))
}

func ConnectMongoDB() {
	uri := getEnv("MONGO_URI", "mongodb://127.0.0.1:27017")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var err error
	MongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Gagal connect MongoDB: %v", err)
	}
	if err = MongoClient.Ping(ctx, nil); err != nil {
		log.Fatalf("Gagal ping MongoDB: %v", err)
	}
	log.Printf("[Config] MongoDB terhubung: %s", uri)
}

func GetMongoCollection(name string) *mongo.Collection {
	return MongoClient.Database(MongoDatabaseName).Collection(name)
}

func CloseConnections() {
	if DB != nil {
		DB.Close()
	}
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		MongoClient.Disconnect(ctx)
	}
}
