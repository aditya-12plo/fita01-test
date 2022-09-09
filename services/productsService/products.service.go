package productsService

import (
	"context"
	"errors"
	"log"
	"os"

	databaseConfig "fita-test-01/config/databaseConfig"
	models "fita-test-01/models"

	"go.mongodb.org/mongo-driver/bson"
)

func ReturnAllProducts() ([]*models.Products, error) {
	var products []*models.Products
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	productsCollection := fita01DB.Collection("products")

	filter := bson.M{}
	cur, err := productsCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cur.Next(context.TODO()) {
		var product models.Products
		err = cur.Decode(&product)
		if err != nil {
			log.Fatal("Error on Decoding the document", err)
		}
		products = append(products, &product)
	}
	defer client.Disconnect(ctx)
	return products, nil

}

func GetProductBySku(skuCode string) (models.Products, error) {
	var product models.Products
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	productsCollection := fita01DB.Collection("products")

	filter := bson.D{{"sku", skuCode}}
	err := productsCollection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return product, err
	}

	return product, nil

}

func UpdateProductBySku(skuCode string, skuQty int) (interface{}, error) {
	var product models.Products
	ctx := context.TODO()

	client := databaseConfig.GetClient()

	fita01DB := client.Database(os.Getenv("DB_DATABASE"))
	productsCollection := fita01DB.Collection("products")

	filter := bson.D{{"sku", skuCode}}
	err := productsCollection.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return product, err
	}

	QtyUpdate := product.Qty - skuQty
	if QtyUpdate < 0 {
		err1 := errors.New("qty not available for sku " + skuCode)
		return nil, err1
	}

	filterUpdate := bson.D{{"_id", product.ID}}
	update := bson.D{{"$set", bson.D{{"qty", QtyUpdate}}}}

	result, errUpdate := productsCollection.UpdateOne(ctx, filterUpdate, update)
	if errUpdate != nil {
		return nil, errUpdate
	}

	return result, nil

}
