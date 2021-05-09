package controller

import (
	"CloudRestaurant/service"
	"CloudRestaurant/tool"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ShopController struct {
}

/**
shop模块的路由解析
*/
func (sc *ShopController) Router(app *gin.Engine) {
	app.GET("/api/shops", sc.GetShopList)
	//测试用 /api/search_shops?keyword=xxx
	app.GET("/api/search_shops", sc.SearchShop)
}

/**
  商铺关信息键词搜索
*/
func (sc *ShopController) SearchShop(context *gin.Context) {
	longitude := context.Query("longitude")
	latitude := context.Query("latitude")
	keyword := context.DefaultQuery("keyword", "")

	if keyword == "" {
		tool.Failed(context, "查询错误,请重新输入商铺名称")
		return
	}
	if longitude == "" || longitude == "undefined" || latitude == "" || latitude == "undefined" {
		longitude = "116.34" //北京
		latitude = "40.34"
	}
	//执行真实的搜索逻辑
	shopService := service.ShopService{}
	shops := shopService.SearchShops(longitude, latitude, keyword)
	if len(shops) != 0 {
		tool.Success(context, shops)
		return
	}
	tool.Failed(context, "未搜索到相关商户")
}

/**
获取商铺列表
*/
func (sc *ShopController) GetShopList(context *gin.Context) {
	longitude := context.Query("longitude")
	latitude := context.Query("latitude")
	//设置默认值
	if longitude == "" || longitude == "undefined" || latitude == "" || latitude == "undefined" {
		longitude = "116.34" //北京
		latitude = "40.34"
	}
	fmt.Println(longitude, latitude)
	shopService := service.ShopService{}
	shops := shopService.ShopList(longitude, latitude)
	if len(shops) == 0 {
		tool.Failed(context, "暂未获取到商户信息")
		return
	}
	//返回之前先用shopid查询其具有的service
	for _, shop := range shops {
		shopServices := shopService.GetService(shop.Id)
		if len(shopServices) == 0 {
			shop.Supports = nil
		} else {
			shop.Supports = shopServices
		}
	}
	tool.Success(context, shops)
}
