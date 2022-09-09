package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Promotions struct {
	ID         primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	PromoCode  string             `bson:"promo_code" json:"promo_code"`
	PromoType  string             `bson:"promo_type" json:"promo_type"`
	Sku        string             `bson:"sku" json:"sku"`
	MinimumQty int                `bson:"minimum_qty" json:"minimum_qty"`
	Discount   float64            `bson:"discount" json:"discount"`
	Details    []detailsPromo     `bson:"details" json:"details_promo"`
}

type detailsPromo struct {
	Sku string `json:"sku"`
	Qty int    `json:"qty"`
}
