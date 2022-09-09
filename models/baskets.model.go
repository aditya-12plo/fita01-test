package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Baskets struct {
	ID      primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	IdBuyer int                `bson:"id_buyer" json:"id_buyer"`
	Sku     string             `json:"sku"`
	Qty     int                `bson:"qty" json:"qty"`
}
