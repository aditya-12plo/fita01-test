package promotionsService

import (
	"context"
	"log"
	"os"

	databaseConfig "fita-test-01/config/databaseConfig"
	models "fita-test-01/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ReturnAllPromotions() ([]*models.Promotions, error) {
	var promotions []*models.Promotions
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	promotionsCollection := fita01DB.Collection("promotions")

	filter := bson.M{}
	cur, err := promotionsCollection.Find(ctx, filter)
	if err != nil {
		log.Fatal("Error on Finding all the documents", err)
	}
	for cur.Next(context.TODO()) {
		var promotion models.Promotions
		err = cur.Decode(&promotion)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		promotions = append(promotions, &promotion)
	}
	defer client.Disconnect(ctx)
	return promotions, nil

}

func GetPromoBySku(skuCode string) (models.Promotions, error) {
	var promo models.Promotions
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	productsCollection := fita01DB.Collection("promotions")

	filter := bson.D{{"sku", skuCode}}
	err := productsCollection.FindOne(ctx, filter).Decode(&promo)
	if err != nil {
		return promo, err
	}

	return promo, nil

}
