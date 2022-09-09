package checkoutsService

import (
	"context"
	"log"
	"os"

	databaseConfig "fita-test-01/config/databaseConfig"
	models "fita-test-01/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ReturnAllCheckouts() ([]*models.Checkouts, error) {
	var checkouts []*models.Checkouts
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	checkoutsCollection := fita01DB.Collection("checkouts")

	filter := bson.M{}
	cur, err := checkoutsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var checkout models.Checkouts
		err = cur.Decode(&checkout)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		checkouts = append(checkouts, &checkout)
	}
	defer client.Disconnect(ctx)
	return checkouts, nil

}

func InsertCheckouts(buyerId int, Sku string, Qty int, Price float64, TotalPrice float64, PromoCode string, PromoType string, Discount float64) (interface{}, error) {

	doc := bson.D{{"id_buyer", buyerId}, {"sku", Sku}, {"qty", Qty}, {"price", Price}, {"total_price", TotalPrice}, {"promo_code", PromoCode}, {"promo_type", PromoType}, {"discount", Discount}}

	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	checkoutsCollection := fita01DB.Collection("checkouts")
	_, err := checkoutsCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
