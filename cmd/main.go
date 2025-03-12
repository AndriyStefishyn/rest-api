package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rest-api/internal/handlers"
	"time"
	"rest-api/internal/shop"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// TODO add context
	clientOptions := options.Client().ApplyURI("mongodb://admin:password@localhost:27017/")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("first-app").Collection("shops")

	MongoStorage := shop.NewMongoStorage(collection)
	MongoStorage.GetShops(ctx)

	router := mux.NewRouter()

	router.HandleFunc("/", handlers.GetAllShopsHandler).Methods(http.MethodGet)
	router.HandleFunc("/getshop/{id}", handlers.GetShopByIdHandler).Methods(http.MethodGet)

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
