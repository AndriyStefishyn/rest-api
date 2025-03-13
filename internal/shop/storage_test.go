package shop

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// remove Global var
// rename file according to go formate style

type StorageTestSuite struct {
	suite.Suite
	ctx            context.Context
	mongoContainer *mongodb.MongoDBContainer
	client         *mongo.Client
	storage        *MongoStorage
	shop           Shop
}

func (suite *StorageTestSuite) SetupSuite() {
	suite.ctx = context.Background()
	mongoContainer, err := mongodb.Run(suite.ctx, "mongo:6")
	// TODO suite.Require().NoError()
	suite.Require().NoError(err)

	suite.mongoContainer = mongoContainer

	uri, err := mongoContainer.ConnectionString(suite.ctx)
	suite.Require().NoError(err)

	client, err := mongo.Connect((suite.ctx), options.Client().ApplyURI(uri))
	suite.Require().NoError(err)

	suite.client = client
	// rename test DB name

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

func (suite *StorageTestSuite) TestUpdateShop() {
	test := Shop{Name: "testshop", Version: 2, Location: "arona", Description: "some test descritpion"}

	err := suite.storage.UpdateShop(suite.ctx, suite.shop.Id, test)
	suite.Require().NoError(err)

	updatedShop, err := suite.storage.GetShopById(suite.ctx, suite.shop.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(updatedShop.Name, test.Name)
	suite.Require().Equal(updatedShop.Version, test.Version)
	suite.Require().Equal(updatedShop.Location, test.Location)
	suite.Require().Equal(updatedShop.Description, test.Description)
}

func (suite *StorageTestSuite) TestGetShopById_Succes() {
	expected, err := suite.storage.GetShopById(suite.ctx, suite.shop.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(expected, suite.shop)
}

func (suite *StorageTestSuite) TestGetShopById_NotFound() {
	suite.shop.Id = "jb5873vf"
	res, _ := suite.storage.GetShopById(suite.ctx, suite.shop.Id)
	suite.Require().Equal(res, Shop{})

}

func (suite *StorageTestSuite) TestGetShops() {
	newshop := Shop{Version: 2, Name: "name2", Location: "location2", Description: "description2"}

	err := suite.storage.InsertShop(suite.ctx, newshop)
	fmt.Println("id", newshop.Id)
	suite.Require().NoError(err)

	shops, err := suite.storage.GetAllShops(suite.ctx)
	fmt.Println(shops)
	suite.Require().NoError(err)

}

func (suite *StorageTestSuite) TestInsertShop() {
	newShop := Shop{Version: 1, Name: "some na", Location: "fjvdknv", Description: "dddfk"}

	err := suite.storage.InsertShop(suite.ctx, newShop)
	suite.Require().NoError(err)

	insertedShop, err := suite.storage.GetShopById(suite.ctx, newShop.Id)
	suite.Require().NoError(err)
	fmt.Println("inserdet shop :", insertedShop.Id)

}
func (suite *StorageTestSuite) TestInsertShop_ID_Exist() {
	testshop := Shop{Version: 1, Name: "some silly shop", Location: "ter", Description: "another silly shop"}

	err := suite.storage.InsertShop(suite.ctx, testshop)
	suite.Require().Error(err)
}

func TestMongoStorageSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
}
