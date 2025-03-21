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
	GetAllShops(ctx context.Context) ([]shop.Shop, error)
}

type ShopService struct {
	storage shop.Storage
}

func NewShopService(storage shop.Storage) *ShopService {
	return &ShopService{storage: storage}
}

func (ss *ShopService) CreateShop(ctx context.Context, shop shop.Shop) (string, error) {
	shop.Id = uuid.NewSHA1(uuid.NameSpaceURL, []byte(shop.Name+shop.Location)).String()

	err := ss.storage.InsertShop(ctx, shop)
	if err != nil {
		return "", fmt.Errorf("create shop:%w", err)
	}

	return shop.Id, nil
}

func (ss *ShopService) DeleteShop(ctx context.Context, id string) error {
	err := ss.storage.DeleteShopById(ctx, id)
	if err != nil {
		return fmt.Errorf("delete shop :%w", err)
	}

	return nil
}

func (ss *ShopService) UpdateShop(ctx context.Context, shop shop.Shop) error {
	err := ss.storage.UpdateShop(ctx, shop)
	if err != nil {
		return fmt.Errorf("update shop:%w", err)
	}

	return nil
}

func (ss *ShopService) GetShop(ctx context.Context, id string) (shop.Shop, error) {
	result, err := ss.storage.GetShopById(ctx, id)
	if err != nil {
		return shop.Shop{}, fmt.Errorf("get shop :%w", err)
	}
	return result, nil
}

func (ss *ShopService) GetAllShops(ctx context.Context) ([]shop.Shop, error) {
	shops, err := ss.storage.GetAllShops(ctx)
	if err != nil {
		return []shop.Shop{}, fmt.Errorf("get all shops :%w", err)
	}
	return shops, nil
}
