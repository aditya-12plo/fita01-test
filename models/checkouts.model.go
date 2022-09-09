package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Checkouts struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	IdBuyer    int                `bson:"id_buyer" json:"id_buyer"`
	Sku        string             `json:"sku"`
	Qty        int                `bson:"qty" json:"qty"`
	Price      float64            `bson:"price" json:"price"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	PromoCode  string             `bson:"promo_code" json:"promo_code"`
	PromoType  string             `bson:"promo_type" json:"promo_type"`
	Discount   float64            `bson:"discount" json:"discount"`
}
