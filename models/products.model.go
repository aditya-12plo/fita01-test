package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Products struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Sku   string             `json:"sku"`
	Name  string             `json:"name"`
	Price float64            `json:"price"`
	Qty   int                `json:"qty"`
}
