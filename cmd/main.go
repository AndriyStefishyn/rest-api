package main

import (
	"context"
	"log"
	"net/http"
	"rest-api/internal/handlers"
	"rest-api/internal/service"
	"rest-api/internal/shop"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	uri := "mongodb://admin:password@localhost:27017/"
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("connection failed")
	}

	collection := client.Database("first-app").Collection("shop")
	storage := shop.NewMongoStorage(collection)
	shopService := service.NewShopService(storage)
	handlers := handlers.NewShopHandlers(shopService)
	router := mux.NewRouter()

	router.HandleFunc("/", handlers.GetShopsHandler).Methods(http.MethodGet)
	router.HandleFunc("/getshop/{id}", handlers.GetShopByIdHandler).Methods(http.MethodGet)
	router.HandleFunc("/createshop", handlers.CreateShopHandler).Methods(http.MethodPost)
	router.HandleFunc("/deleteshop/{id}", handlers.DeleteShopHandler).Methods(http.MethodDelete)
	router.HandleFunc("/updateshop", handlers.UpdateShopHandler).Methods(http.MethodPatch)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal("error server")
	}
}
