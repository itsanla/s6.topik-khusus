package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"nosql-mongodb/internal/config"
	"nosql-mongodb/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	col *mongo.Collection
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{col: config.GetMongoCollection(config.MongoProductCollection)}
}

func newCtx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (r *ProductRepository) Create(p models.Product) error {
	tags, _ := json.Marshal(p.Tags)
	if _, err := config.DB.Exec(
		`INSERT INTO products (product_id,name,price,category,ram,processor,storage,tags,stock)
		 VALUES (?,?,?,?,?,?,?,?,?)`,
		p.ID.Hex(), p.Name, p.Price, p.Category,
		p.Specifications.RAM, p.Specifications.Processor, p.Specifications.Storage,
		string(tags), p.Stock,
	); err != nil {
		return fmt.Errorf("mysql insert product: %w", err)
	}
	ctx, cancel := newCtx()
	defer cancel()
	if _, err := r.col.InsertOne(ctx, p); err != nil {
		return fmt.Errorf("mongo insert product: %w", err)
	}
	return nil
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	ctx, cancel := newCtx()
	defer cancel()
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var products []models.Product
	cursor.All(ctx, &products)
	return products, nil
}

func (r *ProductRepository) GetByID(id primitive.ObjectID) (models.Product, error) {
	ctx, cancel := newCtx()
	defer cancel()
	var p models.Product
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&p); err != nil {
		return p, fmt.Errorf("product not found")
	}
	return p, nil
}

func (r *ProductRepository) Search(filter bson.M) ([]models.Product, error) {
	ctx, cancel := newCtx()
	defer cancel()
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var products []models.Product
	cursor.All(ctx, &products)
	return products, nil
}

func (r *ProductRepository) Update(id string, p models.Product) error {
	tags, _ := json.Marshal(p.Tags)
	if _, err := config.DB.Exec(
		`UPDATE products SET name=?,price=?,category=?,ram=?,processor=?,storage=?,tags=?,stock=? WHERE product_id=?`,
		p.Name, p.Price, p.Category,
		p.Specifications.RAM, p.Specifications.Processor, p.Specifications.Storage,
		string(tags), p.Stock, id,
	); err != nil {
		return fmt.Errorf("mysql update product: %w", err)
	}
	objID, _ := primitive.ObjectIDFromHex(id)
	ctx, cancel := newCtx()
	defer cancel()
	r.col.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": p})
	return nil
}

func (r *ProductRepository) Delete(id string) error {
	config.DB.Exec("DELETE FROM products WHERE product_id=?", id)
	objID, _ := primitive.ObjectIDFromHex(id)
	ctx, cancel := newCtx()
	defer cancel()
	r.col.DeleteOne(ctx, bson.M{"_id": objID})
	return nil
}
