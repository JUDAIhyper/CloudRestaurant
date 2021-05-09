package dao

import (
	"CloudRestaurant/model"
	"CloudRestaurant/tool"
)

type FoodCategoryDao struct {
	*tool.Orm
}

//实例化Dao对象
func NewFoodCategoryDao() *FoodCategoryDao {
	return &FoodCategoryDao{tool.DbEngine}
}

//从数据库中查询所有的食物种类，并返回
func (fcd *FoodCategoryDao) QueryCategories() ([]model.FoodCategory, error) {
	var categories []model.FoodCategory
	err := fcd.DB.Find(&categories).Error
	if err != nil {
		panic(err.Error())
	}
	return categories, nil
}
