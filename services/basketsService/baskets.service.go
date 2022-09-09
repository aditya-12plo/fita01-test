package basketsService

import (
	"context"
	"log"
	"os"

	databaseConfig "fita-test-01/config/databaseConfig"
	models "fita-test-01/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ReturnAllBaskets() ([]*models.Baskets, error) {
	var baskets []*models.Baskets
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	basketsCollection := fita01DB.Collection("baskets")

	filter := bson.M{}
	cur, err := basketsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var basket models.Baskets
		err = cur.Decode(&basket)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		baskets = append(baskets, &basket)
	}
	defer client.Disconnect(ctx)
	return baskets, nil

}

func ReturnAllBasketsBySku(buyerId int, Sku string) ([]*models.Baskets, error) {
	var baskets []*models.Baskets
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	basketsCollection := fita01DB.Collection("baskets")

	filter := bson.D{{"id_buyer", buyerId}, {"sku", Sku}}
	cur, err := basketsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var basket models.Baskets
		err = cur.Decode(&basket)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		baskets = append(baskets, &basket)
	}
	defer client.Disconnect(ctx)
	return baskets, nil

}

func InsertBaskets(buyerId int, Sku string, Qty int) (interface{}, error) {

	doc := bson.D{{"id_buyer", buyerId}, {"sku", Sku}, {"qty", Qty}}

	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	basketsCollection := fita01DB.Collection("baskets")
	_, err := basketsCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}

func DeleteBasketsByBuyyerId(buyerId int) (interface{}, error) {

	doc := bson.D{{"id_buyer", buyerId}}

	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	basketsCollection := fita01DB.Collection("baskets")
	_, err := basketsCollection.DeleteMany(ctx, doc)
	if err != nil {
		return nil, err
	}
	return "ok", nil
}
