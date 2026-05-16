package mongodb

import (
	"context"
	"fmt"

	"northwind-go/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoProductRepository struct {
	col *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) domain.ProductRepository {
	return &mongoProductRepository{col: db.Collection("products")}
}

func (r *mongoProductRepository) FetchActive(ctx context.Context) ([]domain.Product, error) {
	cursor, err := r.col.Find(ctx, bson.M{"discontinued": false})
	if err != nil {
		return nil, fmt.Errorf("mongo find active products: %w", err)
	}
	defer cursor.Close(ctx)
	var products []domain.Product
	if err = cursor.All(ctx, &products); err != nil {
		return nil, fmt.Errorf("mongo decode products: %w", err)
	}
	return products, nil
}

func (r *mongoProductRepository) GetByCode(ctx context.Context, code string) (domain.Product, error) {
	var p domain.Product
	if err := r.col.FindOne(ctx, bson.M{"product_code": code}).Decode(&p); err != nil {
		return p, fmt.Errorf("mongo findone product '%s': %w", code, err)
	}
	return p, nil
}

func (r *mongoProductRepository) Store(ctx context.Context, p *domain.Product) error {
	if _, err := r.col.InsertOne(ctx, p); err != nil {
		return fmt.Errorf("mongo insert product: %w", err)
	}
	return nil
}

func (r *mongoProductRepository) UpdatePrice(ctx context.Context, code string, price float64) error {
	res, err := r.col.UpdateOne(ctx,
		bson.M{"product_code": code},
		bson.M{"$set": bson.M{"list_price": price}},
	)
	if err != nil {
		return fmt.Errorf("mongo update price: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("product '%s' tidak ditemukan", code)
	}
	return nil
}

func (r *mongoProductRepository) SoftDelete(ctx context.Context, code string) error {
	res, err := r.col.UpdateOne(ctx,
		bson.M{"product_code": code},
		bson.M{"$set": bson.M{"discontinued": true}},
	)
	if err != nil {
		return fmt.Errorf("mongo soft delete: %w", err)
	}
	if res.MatchedCount == 0 {
		return fmt.Errorf("product '%s' tidak ditemukan", code)
	}
	return nil
}
