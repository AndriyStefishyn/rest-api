package service

import (
	"context"
	"fmt"
	"rest-api/internal/shop"

	"github.com/google/uuid"
)

type Service interface {
	CreateShop(cxt context.Context, shop shop.Shop) (string, error)
	DeleteShop(cxt context.Context, id string) error
	UpdateShop(cxt context.Context, shop shop.Shop) error
	GetShop(cxt context.Context, id string) (shop.Shop, error)
}

type MongoService struct {
	storage shop.Storage
}

func (ms *MongoService) CreateShop(ctx context.Context, shop shop.Shop) (string, error) {
	shop.Id = fmt.Sprintf("%s,%s,%s", shop.Name, uuid.NewString(), shop.Location)

	err := ms.storage.InsertShop(ctx, shop)
	if err != nil {
		return "", fmt.Errorf("create shop:%w", err)
	}

	return shop.Id, nil
}

func (ms *MongoService) DeleteShop(ctx context.Context, id string) error {
	err := ms.storage.DeleteShopById(ctx, id)
	if err != nil {
		return fmt.Errorf("delete shop :%w", err)
	}

	return nil
}

func (ms *MongoService) UpdateShop(ctx context.Context, shop shop.Shop) error {
	err := ms.storage.UpdateShop(ctx, shop)
	if err != nil {
		return fmt.Errorf("update shop:%w", err)
	}

	return nil
}

func (ms *MongoService) GetShop(ctx context.Context, id string) (shop.Shop, error) {
	result, err := ms.storage.GetShopById(ctx, id)
	if err != nil {
		return shop.Shop{}, fmt.Errorf("get shop :%w", err)
	}
	return result, nil
}
