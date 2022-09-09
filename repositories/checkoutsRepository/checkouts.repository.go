package checkoutsRepository

type DataCheckout struct {
	Sku string `json:"sku"  validate:"required,max=255"`
	Qty int    `json:"qty"  validate:"required,numeric"`
}

type DataCheckoutRepos struct {
	Sku        string
	Qty        int
	CheckSku   interface{}
	CheckPromo interface{}
}

type DetailPromotions struct {
	ID         string          `json:"id"`
	PromoCode  string          `bson:"promo_code" json:"promo_code"`
	PromoType  string          `bson:"promo_type" json:"promo_type"`
	Sku        string          `bson:"sku" json:"sku"`
	MinimumQty int             `bson:"minimum_qty" json:"minimum_qty"`
	Discount   float64         `bson:"discount" json:"discount"`
	Details    []detailsPromos `bson:"details" json:"details_promo"`
}

type detailsPromos struct {
	Sku string `json:"sku"`
	Qty int    `json:"qty"`
}

type DetailProducts struct {
	ID    string  `json:"id"`
	Sku   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Qty   int     `json:"qty"`
}
