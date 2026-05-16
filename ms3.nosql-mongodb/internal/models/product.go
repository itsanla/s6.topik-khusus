package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Specifications struct {
	RAM       string `json:"RAM"       bson:"RAM"`
	Processor string `json:"Processor" bson:"Processor"`
	Storage   string `json:"Storage"   bson:"Storage"`
}

type Product struct {
	ID             primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	Name           string             `json:"name"           bson:"name"`
	Price          int                `json:"price"          bson:"price"`
	Category       string             `json:"category"       bson:"category"`
	Specifications Specifications     `json:"specifications" bson:"specifications"`
	Tags           []string           `json:"tags"           bson:"tags"`
	Stock          int                `json:"stock"          bson:"stock"`
}
