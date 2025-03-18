package shop

import (
	"context"
)

//go:generate mockgen -source=storage.go -destination=mocks/mock_storage.go -package=mocks

type Storage interface {
	GetShopById(ctx context.Context, id string) (Shop, error)
	GetAllShops(ctx context.Context) ([]Shop, error)
	InsertShop(ctx context.Context, shop Shop) error
	UpdateShop(ctx context.Context, shop Shop) error
	DeleteShopById(ctx context.Context, id string) error
}
