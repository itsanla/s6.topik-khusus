package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Customer struct {
	ID            primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	Company       string             `json:"company"        bson:"company"`
	LastName      string             `json:"last_name"      bson:"last_name"`
	FirstName     string             `json:"first_name"     bson:"first_name"`
	EmailAddress  string             `json:"email_address"  bson:"email_address"`
	JobTitle      string             `json:"job_title"      bson:"job_title"`
	BusinessPhone string             `json:"business_phone" bson:"business_phone"`
	City          string             `json:"city"           bson:"city"`
	Country       string             `json:"country_region" bson:"country_region"`
}
