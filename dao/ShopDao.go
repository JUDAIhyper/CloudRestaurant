package dao

import (
	"CloudRestaurant/model"
	"CloudRestaurant/tool"
)

type ShopDao struct {
	*tool.Orm
}

func NewShopDao() *ShopDao {
	return &ShopDao{tool.DbEngine}
}

/**
根据商户id查询对应的服务
*/
func (shopDao *ShopDao) QueryServiceByShopId(shopId int64) []model.Service {
	var service []model.Service
	err := shopDao.Joins("INNER JOIN service on service.id=shop_service.service_id and shop_service.shop_id=?", shopId).Find(&service).Error
	if err != nil {
		return nil
	}
	return service
}

const DEFAULT_RANGE = 5

//116.31  39.21
//115.31-117.31   38.21-40.21

/**
操作数据库查询商铺数据列表
*/
func (shopDao *ShopDao) QueryShops(longitude, latitude float64, keyword string) []model.Shop {
	var shops []model.Shop
	if keyword == "" {
		err := shopDao.DB.Where("longitude> ? and longitude < ? and latitude > ? and latitude < ?",
			longitude-DEFAULT_RANGE, longitude+DEFAULT_RANGE, latitude-DEFAULT_RANGE, latitude+DEFAULT_RANGE).Find(&shops).Error
		if err != nil {
			return nil
		}
	} else {
		err := shopDao.DB.Where("longitude> ? and longitude < ? and latitude > ? and latitude < ? and name like ? and status=1",
			longitude-DEFAULT_RANGE, longitude+DEFAULT_RANGE, latitude-DEFAULT_RANGE, latitude+DEFAULT_RANGE, "%"+keyword+"%").Find(&shops).Error
		if err != nil {
			return nil
		}
	}
	return shops
}
