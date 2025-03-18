package service

import (
	"context"
	"fmt"
	"rest-api/internal/shop"
	"rest-api/internal/shop/mocks"
	"testing"
	"time"

	// indirect

	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"gotest.tools/v3/assert/cmp"
)

type ServiceTestSuite struct {
	suite.Suite
	mockCtrl    *gomock.Controller
	mockStorage *mocks.MockStorage
	service     *ShopService
	ctx         context.Context
	shop        shop.Shop
}

func (suite *ServiceTestSuite) SetupSuite() {
	suite.ctx, _ = context.WithTimeout(context.Background(), time.Second*30)
	suite.mockCtrl = gomock.NewController(suite.T())
	suite.mockStorage = mocks.NewMockStorage(suite.mockCtrl)
	suite.shop = shop.Shop{Id: "1", Version: 1, Name: "shop", Location: "arona", Description: "some fancy shop"}
	suite.service = NewShopService(suite.mockStorage)
}

func (suite *ServiceTestSuite) TearDownSuite() {
	suite.mockCtrl.Finish()
}

func (suite *ServiceTestSuite) TestCreateShop() {
	suite.mockStorage.EXPECT().InsertShop(gomock.Any(), gomock.Any()).DoAndReturn(
		func(ctx context.Context, arg2 shop.Shop) error {
			cmp.DeepEqual(suite.shop, arg2, cmpopts.IgnoreFields(shop.Shop{}, "Id"))

			return nil
		})

	actual, err := suite.service.CreateShop(suite.ctx, suite.shop)
	suite.Require().NoError(err)
	suite.Require().NotEqual("", actual)

}

func (suite *ServiceTestSuite) TestGetShop() {
	newShop := shop.Shop{Id: "1", Version: 1, Name: "shop", Location: "mil", Description: "some fancy shop"}
	suite.mockStorage.EXPECT().GetShopById(gomock.Any(), suite.shop.Id).Return(newShop, nil)

	actual, err := suite.service.GetShop(suite.ctx, "1")
	suite.Require().NoError(err)
	suite.Require().Equal(newShop, actual)

	suite.mockStorage.EXPECT().GetShopById(gomock.Any(), "2").Return(shop.Shop{}, fmt.Errorf("find one"))

	actual, err = suite.mockStorage.GetShopById(suite.ctx, "2")
	suite.Require().Error(err)
	suite.Require().Equal(shop.Shop{}, actual)
}

func (suite *ServiceTestSuite) TestDeleteShop() {
	suite.mockStorage.EXPECT().DeleteShopById(gomock.Any(), suite.shop.Id).Return(nil)

	err := suite.service.DeleteShop(suite.ctx, suite.shop.Id)
	suite.Require().NoError(err)

	suite.mockStorage.EXPECT().DeleteShopById(gomock.Any(), "11").Return(fmt.Errorf("delete"))

	err = suite.service.DeleteShop(suite.ctx, "11")
	suite.Require().Error(err)
}

func (suite *ServiceTestSuite) TestUpdateShop() {
suite.mockStorage.EXPECT().UpdateShop(gomock.Any(),suite.shop).Return(nil)

err := suite.service.UpdateShop(suite.ctx,suite.shop)
suite.Require().NoError(err)

suite.mockStorage.EXPECT().UpdateShop(gomock.Any(),shop.Shop{}).Return(fmt.Errorf("update"))

err = suite.service.storage.UpdateShop(suite.ctx,shop.Shop{})
suite.Require().Error(err)
}
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
