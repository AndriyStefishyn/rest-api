package shop

import (
	"context"
)

type Storage interface {
	// TODO fix methods in interface
	GetShopById(ctx context.Context, id string) (Shop, error)
	GetShops(ctx context.Context) ([]Shop, error)
	InsertShop(ctx context.Context, shop Shop) error
	UpdateShop(ctx context.Context, shop Shop) error
	DeleteShop(ctx context.Context, shop Shop) error
}
