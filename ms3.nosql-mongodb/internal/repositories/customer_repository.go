package repositories

import (
	"context"
	"fmt"
	"time"

	"nosql-mongodb/internal/config"
	"nosql-mongodb/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository struct {
	col *mongo.Collection
}

func NewCustomerRepository() *CustomerRepository {
	return &CustomerRepository{col: config.GetMongoCollection(config.MongoCustomerCollection)}
}

func newCtxC() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

func (r *CustomerRepository) Create(c models.Customer) error {
	if _, err := config.DB.Exec(
		`INSERT INTO customers (customer_id,company,last_name,first_name,email_address,job_title,business_phone,city,country_region)
		 VALUES (?,?,?,?,?,?,?,?,?)`,
		c.ID.Hex(), c.Company, c.LastName, c.FirstName, c.EmailAddress,
		c.JobTitle, c.BusinessPhone, c.City, c.Country,
	); err != nil {
		return fmt.Errorf("mysql insert customer: %w", err)
	}
	ctx, cancel := newCtxC()
	defer cancel()
	if _, err := r.col.InsertOne(ctx, c); err != nil {
		return fmt.Errorf("mongo insert customer: %w", err)
	}
	return nil
}

func (r *CustomerRepository) GetAll() ([]models.Customer, error) {
	ctx, cancel := newCtxC()
	defer cancel()
	cursor, err := r.col.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var customers []models.Customer
	cursor.All(ctx, &customers)
	return customers, nil
}

func (r *CustomerRepository) GetByID(id primitive.ObjectID) (models.Customer, error) {
	ctx, cancel := newCtxC()
	defer cancel()
	var c models.Customer
	if err := r.col.FindOne(ctx, bson.M{"_id": id}).Decode(&c); err != nil {
		return c, fmt.Errorf("customer not found")
	}
	return c, nil
}

func (r *CustomerRepository) Search(filter bson.M) ([]models.Customer, error) {
	ctx, cancel := newCtxC()
	defer cancel()
	cursor, err := r.col.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var customers []models.Customer
	cursor.All(ctx, &customers)
	return customers, nil
}

func (r *CustomerRepository) Delete(id string) error {
	config.DB.Exec("DELETE FROM customers WHERE customer_id=?", id)
	objID, _ := primitive.ObjectIDFromHex(id)
	ctx, cancel := newCtxC()
	defer cancel()
	r.col.DeleteOne(ctx, bson.M{"_id": objID})
	return nil
}
