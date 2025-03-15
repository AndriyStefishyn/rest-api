package shop

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoStorage struct {
	collection *mongo.Collection
}

func NewMongoStorage(collection *mongo.Collection) *MongoStorage {
	return &MongoStorage{collection: collection}
}

func (ms *MongoStorage) GetShopById(ctx context.Context, id string) (Shop, error) {
	var shop Shop

	result := ms.collection.FindOne(ctx, bson.M{"_id": id})
	if result.Err() != nil {
		return Shop{}, fmt.Errorf("find one:%w", result.Err())
	}

	err := result.Decode(&shop)
	if err != nil {
		return Shop{}, fmt.Errorf("decode shop:%w", err)
	}

	return shop, nil
}

func (ms *MongoStorage) GetAllShops(ctx context.Context) ([]Shop, error) {
	cursor, err := ms.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find:%w", err)
	}

	shops, err := cursorToShops(cursor, ctx)
	if err != nil {
		return nil, err
	}

	return shops, nil
}

func cursorToShops(cursor *mongo.Cursor, ctx context.Context) ([]Shop, error) {
	shops := make([]Shop, 0)

	for cursor.Next(ctx) {
		var shop Shop

		if err := cursor.Decode(&shop); err != nil {
			return nil, fmt.Errorf("decode :%w", err)
		}

		shops = append(shops, shop)
	}

	return shops, nil
}

func (ms *MongoStorage) InsertShop(ctx context.Context, shop Shop) error {

	_, err := ms.collection.InsertOne(ctx, shop)

	if err != nil {
		return fmt.Errorf("insert one :%w", err)
	}

	return nil
}

func (ms *MongoStorage) UpdateShop(ctx context.Context, update Shop) error {
	result, err := ms.collection.ReplaceOne(ctx, bson.M{"_id": update.Id}, update)
	if err != nil {
		return fmt.Errorf("update: %w", err)
	}
	
	if result.MatchedCount < 1 {
		return fmt.Errorf("no documents matched")
	}

	return nil
}

func (ms *MongoStorage) DeleteShopById(ctx context.Context, shopId string) error {
	deleteRes, err := ms.collection.DeleteOne(ctx, bson.M{"_id": shopId})
	if err != nil {
		return fmt.Errorf("delete one:%w", err)
	}

	if deleteRes.DeletedCount < 0 {
		return fmt.Errorf("delete one:nothing was deleted")
	}

	return nil
}
