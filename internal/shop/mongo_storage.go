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
		return Shop{}, fmt.Errorf("single result id:%w", result.Err())
	}

	err := result.Decode(&shop)
	if err != nil {
		return Shop{}, fmt.Errorf("decode id:%w", err)
	}

	return shop, nil
}

func (ms *MongoStorage) GetAllShops(ctx context.Context) ([]Shop, error) {
	cursor, err := ms.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("find:%w", err)
	}

	shops, err := convertCursor(cursor, ctx)
	if err != nil {
		return nil, err
	}

	return shops, nil
}

func convertCursor(cursor *mongo.Cursor, ctx context.Context) ([]Shop, error) {
	shops := make([]Shop, 0)

	for cursor.Next(ctx) {
		var shop Shop

		if err := cursor.Decode(&shop); err != nil {
			return nil, fmt.Errorf("decode cursor:%w", err)
		}

		shops = append(shops, shop)
	}

	return shops, nil
}

func (ms *MongoStorage) InsertShop(ctx context.Context, shop Shop) error {

	_, err := ms.collection.InsertOne(ctx, shop)

	if err != nil {
		return fmt.Errorf("insert one %w", err)
	}

	return nil
}

func (ms *MongoStorage) UpdateShop(ctx context.Context, shopId string, update Shop) error {
	filter := bson.M{"_id": shopId}
	res := ms.collection.FindOne(ctx, filter)
	if res.Err() != nil {
		return fmt.Errorf("update find one:%w", res.Err())
	}

	var existingShop Shop
	err := res.Decode(&existingShop)
	if err != nil {
		return fmt.Errorf("update decode :%w", err)
	}

	newShop := Shop{
		Id:          existingShop.Id,
		Version:     update.Version,
		Name:        update.Name,
		Location:    update.Location,
		Description: update.Description,
	}

	result, err := ms.collection.ReplaceOne(ctx, filter, newShop)

	if result.MatchedCount < 1 {
		return fmt.Errorf("nothing to replace")
	}

	if err != nil {
		return fmt.Errorf("update: %w", err)
	}

	return nil
}

func (ms *MongoStorage) DeleteShop(ctx context.Context, shop Shop) error {
	filter := bson.M{"_id": shop.Id}

	deleteRes, err := ms.collection.DeleteOne(ctx, filter)
	if deleteRes.DeletedCount < 0 {
		return fmt.Errorf("delete:nothing was deleted")
	}

	if err != nil {
		return fmt.Errorf("delete:%w", err)
	}

	return nil
}
