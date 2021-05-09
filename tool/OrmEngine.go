package tool

import (
	"CloudRestaurant/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var DbEngine *Orm

type Orm struct {
	*gorm.DB
}

func OrmEngine(cfg *Config) (*Orm, error) {
	database := cfg.Database
	//编写连接的字符串
	conn := database.User + ":" + database.Password + "@tcp(" + database.Host + ":" + database.Port +
		")/" + database.DbName + "?charset=" + database.Charset
	engine, err := gorm.Open(database.Driver, conn)
	if err != nil {
		return nil, err
	}
	//defer engine.Close()
	//全局禁用表名复数
	engine.SingularTable(true)
	//映射数据表
	engine.AutoMigrate(&model.SmsCode{}, &model.Member{}, &model.FoodCategory{},
		&model.Shop{}, &model.Service{}, &model.ShopService{}, &model.Goods{})
	orm := new(Orm)
	orm.DB = engine
	DbEngine = orm

	return orm, nil
}
