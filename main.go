package main

import (
	"CloudRestaurant/controller"
	"CloudRestaurant/middleware"
	"CloudRestaurant/tool"
	"github.com/gin-gonic/gin"
	"log"
)

//路由设置，关联相关服务
func registerRouter(router *gin.Engine) {
	new(controller.HelloController).Router(router)
	new(controller.MemberController).Router(router)
	new(controller.FoodCategoryController).Router(router)
	new(controller.ShopController).Router(router)
	new(controller.GoodController).Router(router)
}

//设置中间件使用
func registerMiddleWare(engine *gin.Engine) {
	//跨域设置
	engine.Use(middleware.Cors())
}

func main() {
	cfg, err := tool.ParseConfig("./config/app.json")
	if err != nil {
		panic(err.Error())
	}
	//调用服务,实例化数据库
	_, err = tool.OrmEngine(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	//初始化redis配置
	tool.InitRedisStore()

	app := gin.Default()
	//使用中间件 *注意：跨域访问必须在配置路由之前使用，否则会失效！
	registerMiddleWare(app)
	//注册路由
	registerRouter(app)
	//集成session功能
	//tool.InitSession(app)
	app.Run(cfg.AppHost + ":" + cfg.AppPort)
}
