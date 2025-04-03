package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"rest-api/internal/handlers"
	"rest-api/internal/service"
	"rest-api/internal/shop"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HandlersTestSuite struct {
	suite.Suite
	mongoContainer *mongodb.MongoDBContainer
	client         *mongo.Client
	context        context.Context
	storage        *shop.MongoStorage
	service        *service.ShopService
	handlers       *handlers.ShopHandlers
	server         *httptest.Server
	router         *mux.Router
	shop           shop.Shop
}

func (suite *HandlersTestSuite) SetupSuite() {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	suite.context = ctx

	mongoContainer, err := mongodb.Run(ctx, "mongo:6")
	suite.Require().NoError(err)

	suite.mongoContainer = mongoContainer

	uri, err := mongoContainer.ConnectionString(ctx)
	suite.Require().NoError(err)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	suite.Require().NoError(err)

	suite.client = client

	suite.shop = shop.Shop{Id: "1", Version: 1, Name: "some name", Location: "here", Description: "this is a fake store"}

	suite.storage = shop.NewMongoStorage(client.Database("testdb").Collection("shops"))
	err = suite.storage.InsertShop(suite.context, suite.shop)
	suite.Require().NoError(err)

	suite.service = service.NewShopService(suite.storage)

	suite.handlers = handlers.NewShopHandlers(suite.service)

	suite.router = mux.NewRouter()
	suite.router.HandleFunc("/", suite.handlers.GetShopsHandler).Methods(http.MethodGet)
	suite.router.HandleFunc("/getshop/{id}", suite.handlers.GetShopByIdHandler).Methods(http.MethodGet)
	suite.router.HandleFunc("/createshop", suite.handlers.CreateShopHandler).Methods(http.MethodPost)
	suite.router.HandleFunc("/deleteshop/{id}", suite.handlers.DeleteShopHandler).Methods(http.MethodDelete)
	suite.router.HandleFunc("/updateshop", suite.handlers.UpdateShopHandler).Methods(http.MethodPut)

	suite.server = httptest.NewServer(suite.router)
}

func (suite *HandlersTestSuite) TearDownSuite() {
	suite.server.Close()
	suite.mongoContainer.Terminate(suite.context)
	suite.client.Disconnect(suite.context)
}

func (suite *HandlersTestSuite) TestGetShopsHandler() {
	resp, err := http.Get(suite.server.URL)
	suite.Require().NoError(err)

	var actual []shop.Shop
	err = json.NewDecoder(resp.Body).Decode(&actual)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.shop, actual[0])

	act, err := http.Get("jdsnvdjndv")
	suite.Require().Error(err)
	suite.Require().Nil(act)
}

func (suite *HandlersTestSuite) TestGetShopByIdHandler() {
	url := fmt.Sprintf("%s/getshop/1", suite.server.URL)

	resp, err := http.Get(url)
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	var actual shop.Shop

	err = json.NewDecoder(resp.Body).Decode(&actual)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.shop, actual)

	url = fmt.Sprintf("%s/getshop/23", suite.server.URL)
	resp, err = http.Get(url)
	suite.Require().NoError(err)
	suite.Require().Empty(resp.Body)
}

func (suite *HandlersTestSuite) TestCreateShopHandler() {
	url := fmt.Sprintf("%s/createshop", suite.server.URL)

	testShop := shop.Shop{Version: 1, Name: "some name", Location: "here", Description: "this is a fake store"}

	jsonData, err := json.Marshal(testShop)
	suite.Require().NoError(err)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	suite.Require().Equal(http.StatusOK, resp.StatusCode)

	var id string
	err = json.NewDecoder(resp.Body).Decode(&id)
	suite.Require().NoError(err)

	url = fmt.Sprintf("%s/getshop/%s", suite.server.URL, id)

	resp, err = http.Get(url)
	suite.Require().NoError(err)

	var expected shop.Shop

	err = json.NewDecoder(resp.Body).Decode(&expected)
	suite.Require().NoError(err)

	expected = shop.Shop{Version: expected.Version, Name: expected.Name, Location: expected.Location, Description: expected.Description}
	suite.Require().Equal(expected, testShop)
}

func (suite *HandlersTestSuite) TestDeleteShopHandler() {
	url := fmt.Sprintf("%s/deleteshop/%s", suite.server.URL, suite.shop.Id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	suite.Require().NoError(err)

	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)
	suite.Require().NotEmpty(rr.Body)

	resp, err := http.Get(url)
	suite.Require().NoError(err)
	suite.Require().Empty(resp.Body)
}

func (suite *HandlersTestSuite) TestUpdateShopHandler() {
	testShop := shop.Shop{Id: "1", Version: 1, Name: "some updates", Location: "update", Description: "updated shop"}
	jsonData, err := json.Marshal(testShop)
	suite.Require().NoError(err)

	url := fmt.Sprintf("%s/updateshop", suite.server.URL)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	rr := httptest.NewRecorder()

	suite.router.ServeHTTP(rr, req)

	suite.Require().NotEmpty(rr.Body)

	url = fmt.Sprintf("%s/getshop/%s", suite.server.URL, suite.shop.Id)
	resp, err := http.Get(url)
	suite.Require().NoError(err)

	var expected shop.Shop

	err = json.NewDecoder(resp.Body).Decode(&expected)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, testShop)
}
func TestHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(HandlersTestSuite))
}
