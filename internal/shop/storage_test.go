package shop

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StorageTestSuite struct {
	suite.Suite
	ctx            context.Context
	mongoContainer *mongodb.MongoDBContainer
	client         *mongo.Client
	storage        *MongoStorage
	shop           Shop
}

func (suite *StorageTestSuite) SetupSuite() {
	suite.ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)

	mongoContainer, err := mongodb.Run(suite.ctx, "mongo:6")
	suite.Require().NoError(err)

	suite.mongoContainer = mongoContainer

	uri, err := mongoContainer.ConnectionString(suite.ctx)
	suite.Require().NoError(err)

	client, err := mongo.Connect((suite.ctx), options.Client().ApplyURI(uri))
	suite.Require().NoError(err)

	suite.client = client

	suite.storage = NewMongoStorage(client.Database("first-app").Collection("shops"))

	newId := uuid.New().String()

	suite.shop = Shop{Id: newId, Version: 1, Name: "name", Location: "location", Description: "description"}

	err = suite.storage.InsertShop(suite.ctx, suite.shop)
	suite.Require().NoError(err)
}

func (suite *StorageTestSuite) TearDownSuite() {
	err := suite.client.Disconnect(suite.ctx)
	suite.Require().NoError(err)

	err = suite.mongoContainer.Terminate(suite.ctx)
	suite.Require().NoError(err)
}

func (suite *StorageTestSuite) TestGetShopById_Succes() {
	expected := suite.shop

	actual, err := suite.storage.GetShopById(suite.ctx, suite.shop.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, actual)
}

func (suite *StorageTestSuite) TestGetShopById_NotFound() {
	suite.shop.Id = ""

	res, err := suite.storage.GetShopById(suite.ctx, suite.shop.Id)
	suite.Require().Error(err)
	suite.Require().Equal(res, Shop{})
}

func (suite *StorageTestSuite) TestGetAllShops() {
	newshop := Shop{Version: 2, Name: "name2", Location: "location2", Description: "description2"}

	err := suite.storage.InsertShop(suite.ctx, newshop)
	suite.Require().NoError(err)

	expected, err := suite.storage.GetAllShops(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(suite.shop, expected[0])
	suite.Require().Equal(newshop, expected[1])
}

func (suite *StorageTestSuite) TestInsertShop() {
	newShop := Shop{Id: "2", Version: 1, Name: "some name", Location: "fjvdknv", Description: "dddfk"}

	err := suite.storage.InsertShop(suite.ctx, newShop)
	suite.Require().NoError(err)

	actual, err := suite.storage.GetShopById(suite.ctx, newShop.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(newShop, actual)
}

func (suite *StorageTestSuite) TestUpdateShop() {

	update := Shop{Id: suite.shop.Id, Name: "testshop", Version: 2, Location: "arona", Description: "some test descritpion"}

	err := suite.storage.UpdateShop(suite.ctx, update)
	suite.Require().NoError(err)

	actual, err := suite.storage.GetShopById(suite.ctx, suite.shop.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(update,actual)
}

func (suite *StorageTestSuite) TestUpdateShop_No_File() {
	test := Shop{Name: "testshop", Version: 2, Location: "arona", Description: "some test descritpion"}

	err := suite.storage.UpdateShop(suite.ctx, test)
	suite.Require().Error(err)
}

func (suite *StorageTestSuite) TestDeleteShop() {
	shopIdToDelete := suite.shop.Id

	err := suite.storage.DeleteShopById(suite.ctx, shopIdToDelete)
	suite.Require().NoError(err)

	actual, err := suite.storage.GetShopById(suite.ctx, shopIdToDelete)
	suite.Require().Error(err)
	suite.Require().Equal(Shop{}, actual)
}

func (suite *StorageTestSuite) TestDeleteShop_Failed() {
	wrongId := "some37845vfe"

	err := suite.storage.DeleteShopById(suite.ctx, wrongId)
	suite.Require().Error(err)
}

func TestMongoStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
