package shop

import (
	"context"
)

type Storage interface {
	GetShopById(ctx context.Context, id string) (Shop, error)
	GetAllShops(ctx context.Context) ([]Shop, error)
	InsertShop(ctx context.Context, shop Shop) error
	UpdateShop(ctx context.Context, shop Shop) error
	DeleteShopById(ctx context.Context, id string) error
}
